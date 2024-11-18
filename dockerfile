FROM mcr.microsoft.com/microsoftgraph/powershell:latest

WORKDIR /app
COPY scripts/ExportGraphPermissions.ps1 /app/

ENTRYPOINT ["pwsh", "-File", "/app/ExportGraphPermissions.ps1"]