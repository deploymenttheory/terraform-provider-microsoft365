#!/usr/bin/env python3
"""
Generate proof-of-possession (PoP) tokens for automating key rolling
and certificate updates using Microsoft Graph.

This script generates:
1. Client assertion tokens for authentication
2. Proof-of-possession (PoP) tokens for addKey/removeKey operations

Based on Microsoft documentation:
https://learn.microsoft.com/en-us/graph/application-rollkey-prooftoken

The PoP token must contain:
- aud: Audience must be 00000002-0000-0000-c000-000000000000
- iss: Issuer should be the application or service principal ID
- nbf: Not before time
- exp: Expiration time (nbf + 10 minutes recommended)

Usage:
    python generate_pop_token.py --client-id <app_id> --object-id <object_id> --cert-path <path> [--password <pwd>]

Requirements:
    pip install cryptography PyJWT

Author: Terraform Provider Microsoft365 Team
"""

import argparse
import base64
import json
import sys
import time
import uuid
from datetime import datetime, timezone
from pathlib import Path
from typing import Optional, Tuple

try:
    import jwt
    from cryptography import x509
    from cryptography.hazmat.backends import default_backend
    from cryptography.hazmat.primitives import hashes, serialization
    from cryptography.hazmat.primitives.serialization import pkcs12
except ImportError as e:
    print(f"Error: Missing required dependency: {e}")
    print("Please install required packages: pip install cryptography PyJWT")
    sys.exit(1)


# Constants
POP_TOKEN_AUDIENCE = "00000002-0000-0000-c000-000000000000"
TOKEN_LIFETIME_MINUTES = 10


def load_certificate(
    cert_path: str, password: Optional[str] = None
) -> Tuple[x509.Certificate, any]:
    """
    Load a certificate and private key from a PFX/P12 or PEM file.

    Args:
        cert_path: Path to the certificate file (.pfx, .p12, or .pem)
        password: Optional password for the certificate

    Returns:
        Tuple of (certificate, private_key)

    Raises:
        ValueError: If the certificate format is not supported
        FileNotFoundError: If the certificate file doesn't exist
    """
    cert_path = Path(cert_path)

    if not cert_path.exists():
        raise FileNotFoundError(f"Certificate file not found: {cert_path}")

    with open(cert_path, "rb") as f:
        cert_data = f.read()

    password_bytes = password.encode() if password else None

    # Try to load as PFX/P12 first
    if cert_path.suffix.lower() in (".pfx", ".p12"):
        private_key, certificate, _ = pkcs12.load_key_and_certificates(
            cert_data, password_bytes, default_backend()
        )
        return certificate, private_key

    # Try to load as PEM
    if cert_path.suffix.lower() in (".pem", ".crt", ".cer"):
        # Load certificate
        certificate = x509.load_pem_x509_certificate(cert_data, default_backend())

        # Try to load private key from the same file or separate .key file
        try:
            private_key = serialization.load_pem_private_key(
                cert_data, password=password_bytes, backend=default_backend()
            )
        except ValueError as exc:
            # Try loading from a separate .key file
            key_path = cert_path.with_suffix(".key")
            if key_path.exists():
                with open(key_path, "rb") as f:
                    key_data = f.read()
                private_key = serialization.load_pem_private_key(
                    key_data, password=password_bytes, backend=default_backend()
                )
            else:
                raise ValueError(
                    f"Private key not found in {cert_path} or {key_path}"
                ) from exc

        return certificate, private_key

    raise ValueError(f"Unsupported certificate format: {cert_path.suffix}")


def get_certificate_thumbprint(certificate: x509.Certificate) -> str:
    """
    Get the SHA-1 thumbprint of a certificate (base64url encoded).

    Args:
        certificate: The X.509 certificate

    Returns:
        Base64url encoded SHA-1 thumbprint
    """
    thumbprint = certificate.fingerprint(hashes.SHA1())
    return base64.urlsafe_b64encode(thumbprint).decode().rstrip("=")


def get_certificate_x5c(certificate: x509.Certificate) -> str:
    """
    Get the X.509 certificate chain (x5c) value for JWT header.

    Args:
        certificate: The X.509 certificate

    Returns:
        Base64 encoded certificate (standard encoding, not URL-safe)
    """
    cert_der = certificate.public_bytes(serialization.Encoding.DER)
    return base64.b64encode(cert_der).decode()


def get_certificate_public_key_base64(certificate: x509.Certificate) -> str:
    """
    Get the base64-encoded public key from a certificate.
    This is the value used for the 'key' field in addKey requests.

    Args:
        certificate: The X.509 certificate

    Returns:
        Base64 encoded certificate in DER format
    """
    cert_der = certificate.public_bytes(serialization.Encoding.DER)
    return base64.b64encode(cert_der).decode()


def generate_pop_token(
    object_id: str,
    certificate: x509.Certificate,
    private_key: any,
    lifetime_minutes: int = TOKEN_LIFETIME_MINUTES,
) -> str:
    """
    Generate a proof-of-possession (PoP) token for addKey/removeKey operations.

    The PoP token is a self-signed JWT with the following claims:
    - aud: 00000002-0000-0000-c000-000000000000
    - iss: The object ID of the application or service principal
    - nbf: Current time
    - exp: Current time + lifetime_minutes

    Args:
        object_id: The object ID of the application or service principal
        certificate: The X.509 certificate to sign with
        private_key: The private key to sign with
        lifetime_minutes: Token lifetime in minutes (default: 10)

    Returns:
        The signed JWT PoP token
    """
    current_time = int(time.time())
    expiration_time = current_time + (lifetime_minutes * 60)

    # JWT header with x5t (thumbprint)
    thumbprint = get_certificate_thumbprint(certificate)

    headers = {
        "alg": "RS256",
        "typ": "JWT",
        "x5t": thumbprint,
    }

    # JWT payload
    payload = {
        "aud": POP_TOKEN_AUDIENCE,
        "iss": object_id,
        "nbf": current_time,
        "exp": expiration_time,
    }

    # Sign the token
    token = jwt.encode(payload, private_key, algorithm="RS256", headers=headers)

    return token


def generate_client_assertion(
    client_id: str,
    tenant_id: str,
    certificate: x509.Certificate,
    private_key: any,
    lifetime_minutes: int = TOKEN_LIFETIME_MINUTES,
) -> str:
    """
    Generate a client assertion token for authenticating with Azure AD.

    This token is used to obtain an access token using client credentials flow
    with certificate authentication.

    Args:
        client_id: The application (client) ID
        tenant_id: The Azure AD tenant ID
        certificate: The X.509 certificate to sign with
        private_key: The private key to sign with
        lifetime_minutes: Token lifetime in minutes (default: 10)

    Returns:
        The signed JWT client assertion token
    """
    current_time = int(time.time())
    expiration_time = current_time + (lifetime_minutes * 60)

    # Audience for client assertion
    audience = f"https://login.microsoftonline.com/{tenant_id}/v2.0"

    # JWT header with x5t (thumbprint) and x5c (certificate chain)
    thumbprint = get_certificate_thumbprint(certificate)
    x5c = get_certificate_x5c(certificate)

    headers = {
        "alg": "RS256",
        "typ": "JWT",
        "x5t": thumbprint,
        "x5c": [x5c],
    }

    # JWT payload
    payload = {
        "aud": audience,
        "iss": client_id,
        "sub": client_id,
        "jti": str(uuid.uuid4()),
        "nbf": current_time,
        "exp": expiration_time,
    }

    # Sign the token
    token = jwt.encode(payload, private_key, algorithm="RS256", headers=headers)

    return token


def decode_token_for_display(token: str) -> dict:
    """
    Decode a JWT token without verification for display purposes.

    Args:
        token: The JWT token to decode

    Returns:
        Dictionary containing header and payload
    """
    try:
        # Decode without verification
        header = jwt.get_unverified_header(token)
        payload = jwt.decode(token, options={"verify_signature": False})

        # Convert timestamps to readable format
        if "nbf" in payload:
            payload["nbf_readable"] = datetime.fromtimestamp(
                payload["nbf"], tz=timezone.utc
            ).isoformat()
        if "exp" in payload:
            payload["exp_readable"] = datetime.fromtimestamp(
                payload["exp"], tz=timezone.utc
            ).isoformat()

        return {"header": header, "payload": payload}
    except (jwt.exceptions.DecodeError, KeyError) as e:
        return {"error": str(e)}


def print_token_info(token_name: str, token: str):
    """
    Print token information in a formatted way.

    Args:
        token_name: Name/description of the token
        token: The JWT token
    """
    print(f"\n{'=' * 60}")
    print(f"{token_name}")
    print("=" * 60)
    print(f"\nToken Value:\n{token}")

    decoded = decode_token_for_display(token)
    if "error" not in decoded:
        print(f"\nDecoded Header:\n{json.dumps(decoded['header'], indent=2)}")
        print(f"\nDecoded Payload:\n{json.dumps(decoded['payload'], indent=2)}")
    print()


def main():
    """Entry point for the PoP token generator CLI."""
    parser = argparse.ArgumentParser(
        description="Generate proof-of-possession (PoP) tokens for Microsoft Graph key operations",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Generate PoP token for addKey/removeKey operations
  python generate_pop_token.py --object-id <object_id> --cert-path ./cert.pfx --password secret

  # Generate both client assertion and PoP token
  python generate_pop_token.py --client-id <app_id> --tenant-id <tenant_id> --object-id <object_id> --cert-path ./cert.pfx

  # Display certificate information
  python generate_pop_token.py --cert-path ./cert.pfx --info-only

For more information, see:
https://learn.microsoft.com/en-us/graph/application-rollkey-prooftoken
        """,
    )

    parser.add_argument(
        "--client-id",
        help="Application (client) ID - required for client assertion",
    )
    parser.add_argument(
        "--tenant-id",
        help="Azure AD tenant ID - required for client assertion",
    )
    parser.add_argument(
        "--object-id",
        help="Object ID of the application or service principal - required for PoP token",
    )
    parser.add_argument(
        "--cert-path",
        required=True,
        help="Path to the certificate file (.pfx, .p12, or .pem)",
    )
    parser.add_argument(
        "--password",
        help="Password for the certificate (if required)",
    )
    parser.add_argument(
        "--lifetime",
        type=int,
        default=TOKEN_LIFETIME_MINUTES,
        help=f"Token lifetime in minutes (default: {TOKEN_LIFETIME_MINUTES})",
    )
    parser.add_argument(
        "--output-format",
        choices=["text", "json"],
        default="text",
        help="Output format (default: text)",
    )
    parser.add_argument(
        "--info-only",
        action="store_true",
        help="Only display certificate information, don't generate tokens",
    )
    parser.add_argument(
        "--new-cert-path",
        help="Path to a new certificate to get the key value for addKey requests",
    )
    parser.add_argument(
        "--new-cert-password",
        help="Password for the new certificate (if required)",
    )

    args = parser.parse_args()

    try:
        # Load the existing certificate
        print(f"\nLoading certificate from: {args.cert_path}")
        certificate, private_key = load_certificate(args.cert_path, args.password)

        # Display certificate information
        thumbprint = get_certificate_thumbprint(certificate)
        print("Certificate loaded successfully")
        print(f"  Subject: {certificate.subject.rfc4514_string()}")
        print(f"  Issuer: {certificate.issuer.rfc4514_string()}")
        print(f"  Valid From: {certificate.not_valid_before_utc}")
        print(f"  Valid Until: {certificate.not_valid_after_utc}")
        print(f"  Thumbprint (x5t): {thumbprint}")

        if args.info_only:
            # If a new certificate is specified, show its key value
            if args.new_cert_path:
                print(f"\nLoading new certificate from: {args.new_cert_path}")
                new_cert, _ = load_certificate(
                    args.new_cert_path, args.new_cert_password
                )
                key_value = get_certificate_public_key_base64(new_cert)
                print("\nNew certificate key value (for addKey 'key' field):")
                print(key_value)
            return

        result = {}

        # Generate client assertion if tenant_id and client_id are provided
        if args.client_id and args.tenant_id:
            client_assertion = generate_client_assertion(
                args.client_id,
                args.tenant_id,
                certificate,
                private_key,
                args.lifetime,
            )
            result["client_assertion"] = client_assertion

            if args.output_format == "text":
                print_token_info("Client Assertion Token", client_assertion)

        # Generate PoP token if object_id is provided
        if args.object_id:
            pop_token = generate_pop_token(
                args.object_id,
                certificate,
                private_key,
                args.lifetime,
            )
            result["pop_token"] = pop_token
            result["proof"] = pop_token  # Alias for convenience

            if args.output_format == "text":
                print_token_info("Proof-of-Possession (PoP) Token", pop_token)

        # If a new certificate is specified, get its key value
        if args.new_cert_path:
            print(f"\nLoading new certificate from: {args.new_cert_path}")
            new_cert, _ = load_certificate(args.new_cert_path, args.new_cert_password)
            key_value = get_certificate_public_key_base64(new_cert)
            result["new_certificate_key"] = key_value

            if args.output_format == "text":
                print(f"\n{'=' * 60}")
                print("New Certificate Key Value (for addKey 'key' field)")
                print("=" * 60)
                print(f"\n{key_value}\n")

        # Output JSON if requested
        if args.output_format == "json":
            print(json.dumps(result, indent=2))

        # Show usage hints
        if args.output_format == "text" and (result.get("pop_token") or result.get("client_assertion")):
            print("\n" + "=" * 60)
            print("Usage Hints")
            print("=" * 60)

            if result.get("pop_token"):
                print("\nUse the PoP token as the 'proof' parameter in addKey/removeKey requests:")
                print("""
    POST https://graph.microsoft.com/v1.0/applications/{object-id}/addKey
    Content-Type: application/json

    {
        "keyCredential": {
            "type": "AsymmetricX509Cert",
            "usage": "Verify",
            "key": "<base64-encoded-certificate>"
        },
        "passwordCredential": null,
        "proof": "<pop_token>"
    }
""")

            if result.get("client_assertion"):
                print("\nUse the client assertion to obtain an access token:")
                print(f"""
    POST https://login.microsoftonline.com/{args.tenant_id}/oauth2/v2.0/token
    Content-Type: application/x-www-form-urlencoded

    client_id={args.client_id}
    &scope=https://graph.microsoft.com/.default
    &client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer
    &client_assertion=<client_assertion>
    &grant_type=client_credentials
""")

    except FileNotFoundError as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
    except ValueError as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
    except (OSError, TypeError) as e:
        print(f"Unexpected error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()

