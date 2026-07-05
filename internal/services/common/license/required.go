package license

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// RequiredLicenses defines the licenses required for specific features
// Each feature lists all possible licenses (SKUs or service plans) that enable it
//
// Reference: https://m365maps.com/ provides comprehensive Microsoft 365 licensing diagrams
// and shows which features are included in each license tier.
var RequiredLicenses = map[string][]string{
	// ========================================================================
	// Microsoft Entra (Identity & Access Management)
	// Reference: https://m365maps.com/ - Entra ID section
	// ========================================================================

	"ConditionalAccessPolicy": {
		constants.SKUMicrosoftEntraSuite,  // Entra Suite (includes P2)
		constants.SKUSPEE5,                // Microsoft 365 E5 (includes Entra ID P1)
		constants.ServicePlanAADPREMIUM,   // Entra ID P1 service plan
		constants.ServicePlanAADPREMIUMP2, // Entra ID P2 service plan
	},

	"PrivilegedIdentityManagement": {
		constants.SKUMicrosoftEntraSuite,   // Entra Suite (includes P2)
		constants.SKUSPEE5,                 // Microsoft 365 E5 (includes Entra ID P2)
		constants.ServicePlanAADPREMIUMP2,  // Entra ID P2 (includes PIM)
		constants.ServicePlanPAMENTERPRISE, // PAM Enterprise service plan
	},

	"IdentityGovernance": {
		constants.SKUMicrosoftEntraSuite,             // Entra Suite
		constants.SKUSPEE5,                           // Microsoft 365 E5
		constants.ServicePlanEntraIDentityGovernance, // Identity Governance service plan
	},

	"VerifiedID": {
		constants.SKUMicrosoftEntraSuite,                         // Entra Suite
		constants.ServicePlanVerifiableCredentialsServiceRequest, // Verifiable Credentials service plan
	},

	"NetworkFilteringPolicy": {
		constants.SKUMicrosoftEntraSuite,                // Microsoft Entra Suite (includes Global Secure Access)
		constants.ServicePlanEntraPremiumInternetAccess, // Internet Access service plan
	},

	"PrivateAccessPolicy": {
		constants.SKUMicrosoftEntraSuite,               // Microsoft Entra Suite (includes Global Secure Access)
		constants.ServicePlanEntraPremiumPrivateAccess, // Private Access service plan
	},

	"MultiFactorAuthentication": {
		constants.SKUMicrosoftEntraSuite, // Entra Suite
		constants.SKUSPEE5,               // Microsoft 365 E5
		constants.ServicePlanAADPREMIUM,  // Entra ID P1 (includes MFA)
		constants.ServicePlanMFAPREMIUM,  // MFA Premium service plan
	},

	// ========================================================================
	// Microsoft Intune (Endpoint Management)
	// Reference: https://m365maps.com/ - Intune section
	// ========================================================================

	"IntuneBasicDeviceManagement": {
		constants.SKUMicrosoftIntuneSuite, // Microsoft Intune Suite
		constants.SKUSPEE5,                // Microsoft 365 E5 (includes Intune)
		constants.ServicePlanINTUNEA,      // Intune standalone service plan
		constants.ServicePlanINTUNEO365,   // Intune for Office 365 service plan
	},

	"IntuneEndpointPrivilegeManagement": {
		constants.SKUMicrosoftIntuneSuite, // Microsoft Intune Suite
		constants.ServicePlanIntuneEPM,    // Intune Endpoint Privilege Management service plan
	},

	"IntuneCloudPKI": {
		constants.SKUMicrosoftIntuneSuite, // Microsoft Intune Suite
		constants.ServicePlanCLOUDPKI,     // Cloud PKI service plan
	},

	"IntuneAdvancedEndpointAnalytics": {
		constants.SKUMicrosoftIntuneSuite,     // Microsoft Intune Suite
		constants.ServicePlanIntuneAdvancedea, // Advanced Endpoint Analytics service plan
	},

	"IntuneRemoteHelp": {
		constants.SKUMicrosoftIntuneSuite, // Microsoft Intune Suite
		constants.ServicePlanREMOTEHELP,   // Remote Help service plan
	},

	"IntuneMAMTunnel": {
		constants.SKUMicrosoftIntuneSuite,    // Microsoft Intune Suite
		constants.ServicePlanIntuneMamtunnel, // MAM Tunnel service plan
	},

	"IntuneServiceNowIntegration": {
		constants.SKUMicrosoftIntuneSuite,     // Microsoft Intune Suite
		constants.ServicePlanIntuneServicenow, // ServiceNow integration service plan
	},

	"IntuneThirdPartyAppPatching": {
		constants.SKUMicrosoftIntuneSuite,   // Microsoft Intune Suite
		constants.ServicePlan3PARTYAPPPATCH, // Third-party app patching service plan
	},

	"WindowsAutopatch": {
		constants.SKUSPEE5,                    // Microsoft 365 E5
		constants.ServicePlanWindowsAutopatch, // Windows Autopatch service plan
	},

	// ========================================================================
	// Windows 365 (Cloud PC)
	// Reference: https://m365maps.com/ - Windows 365 section
	// ========================================================================

	"Windows365CloudPC": {
		constants.SKUCPCB2C8RAM128GB,          // Windows 365 Business
		constants.SKUCPCE2C8GB128GB,           // Windows 365 Enterprise
		constants.SKUWindows365S2vCPU8GB128GB, // Windows 365 Frontline
		constants.ServicePlanCPC2,             // Cloud PC E service plan
		constants.ServicePlanCPCS2C8GB128GB,   // Cloud PC S service plan
		constants.ServicePlanCPCSS2,           // Cloud PC Shared Services
	},

	"Windows10ESU": {
		constants.SKUCPCB2C8RAM128GB,                // Windows 365 Business (includes ESU)
		constants.SKUCPCE2C8GB128GB,                 // Windows 365 Enterprise (includes ESU)
		constants.ServicePlanWindows10ESUCommercial, // Windows 10 ESU Commercial service plan
		constants.ServicePlanWINDOWS10ESUTENANT,     // Windows 10 ESU Tenant service plan
	},

	// ========================================================================
	// Microsoft Defender (Threat Protection)
	// Reference: https://m365maps.com/ - Defender section
	// ========================================================================

	"DefenderForEndpoint": {
		constants.SKUSPEE5,             // Microsoft 365 E5
		constants.ServicePlanWINDEFATP, // Windows Defender ATP service plan
	},

	"DefenderForOffice365": {
		constants.SKUSPEE5,                                   // Microsoft 365 E5
		constants.ServicePlanATPENTERPRISE,                   // ATP Enterprise service plan
		constants.ServicePlanCOMMONDEFENDERPLATFORMFOROFFICE, // Common Defender Platform for Office
	},

	"DefenderForCloudApps": {
		constants.SKUSPEE5,                      // Microsoft 365 E5
		constants.ServicePlanADALLOMSO365,       // Defender for Cloud Apps for O365
		constants.ServicePlanADALLOMSSTANDALONE, // Defender for Cloud Apps Standalone
	},

	"DefenderForIdentity": {
		constants.SKUSPEE5,       // Microsoft 365 E5
		constants.ServicePlanATA, // Advanced Threat Analytics (Defender for Identity)
	},

	"DefenderXDR": {
		constants.SKUSPEE5,       // Microsoft 365 E5
		constants.ServicePlanMTP, // Microsoft 365 Defender (XDR)
	},

	"DefenderThreatIntelligence": {
		constants.SKUSPEE5,                      // Microsoft 365 E5
		constants.ServicePlanTHREATINTELLIGENCE, // Threat Intelligence service plan
	},

	"DefenderForIoT": {
		constants.SKUSPEE5, // Microsoft 365 E5
		constants.ServicePlanDefenderForIotEnterprise, // Defender for IoT Enterprise service plan
	},

	// ========================================================================
	// Microsoft Purview (Compliance & Data Governance)
	// Reference: https://m365maps.com/ - Purview section
	// ========================================================================

	"PurviewDataLifecycleManagement": {
		constants.SKUSPEE5,                     // Microsoft 365 E5
		constants.ServicePlanINFOGOVERNANCE,    // Information Governance service plan
		constants.ServicePlanRECORDSMANAGEMENT, // Records Management service plan
	},

	"PurviewInformationProtection": {
		constants.SKUSPEE5,                  // Microsoft 365 E5
		constants.ServicePlanMIPSCLP1,       // MIP Sensitivity Labels P1
		constants.ServicePlanMIPSCLP2,       // MIP Sensitivity Labels P2
		constants.ServicePlanMIPSExchange,   // MIP for Exchange
		constants.ServicePlanRMSSENTERPRISE, // Rights Management Service Enterprise
		constants.ServicePlanRMSSPREMIUM,    // Rights Management Service Premium
		constants.ServicePlanRMSSPREMIUM2,   // Rights Management Service Premium 2
	},

	"PurviewDataLossPrevention": {
		constants.SKUSPEE5,                        // Microsoft 365 E5
		constants.ServicePlanCOMMUNICATIONSDLP,    // Communications DLP service plan
		constants.ServicePlanMICROSOFTENDPOINTDLP, // Endpoint DLP service plan
	},

	"PurviewInsiderRiskManagement": {
		constants.SKUSPEE5,                         // Microsoft 365 E5
		constants.ServicePlanINSIDERRISK,           // Insider Risk Management service plan
		constants.ServicePlanINSIDERRISKMANAGEMENT, // Insider Risk Management (full) service plan
	},

	"PurviewCommunicationCompliance": {
		constants.SKUSPEE5, // Microsoft 365 E5
		constants.ServicePlanCOMMUNICATIONSCOMPLIANCE,         // Communications Compliance service plan
		constants.ServicePlanMICROSOFTCOMMUNICATIONCOMPLIANCE, // Microsoft Communication Compliance
	},

	"PurvieweDiscoveryPremium": {
		constants.SKUSPEE5,                    // Microsoft 365 E5
		constants.ServicePlanEQUIVIOANALYTICS, // eDiscovery (Equivio) Analytics service plan
	},

	"PurviewAdvancedAudit": {
		constants.SKUSPEE5,                        // Microsoft 365 E5
		constants.ServicePlanM365ADVANCEDAUDITING, // M365 Advanced Auditing service plan
		constants.ServicePlanM365AUDITPLATFORM,    // M365 Audit Platform service plan
	},

	"PurviewContentExplorer": {
		constants.SKUSPEE5,                           // Microsoft 365 E5
		constants.ServicePlanContentExplorer,         // Content Explorer service plan
		constants.ServicePlanContentexplorerStandard, // Content Explorer Standard service plan
	},

	"PurviewDataInvestigations": {
		constants.SKUSPEE5,                      // Microsoft 365 E5
		constants.ServicePlanDATAINVESTIGATIONS, // Data Investigations service plan
	},

	"PurviewInformationBarriers": {
		constants.SKUSPEE5,                       // Microsoft 365 E5
		constants.ServicePlanINFORMATIONBARRIERS, // Information Barriers service plan
	},

	"PurviewCustomerKey": {
		constants.SKUSPEE5,               // Microsoft 365 E5
		constants.ServicePlanCUSTOMERKEY, // Customer Key service plan
	},

	"PurviewCustomerLockbox": {
		constants.SKUSPEE5, // Microsoft 365 E5
		constants.ServicePlanCustomerlockboxaEnterprise, // Customer Lockbox Enterprise service plan
		constants.ServicePlanLOCKBOXENTERPRISE,          // Lockbox Enterprise service plan
	},

	"PurviewPremiumEncryption": {
		constants.SKUSPEE5,                     // Microsoft 365 E5
		constants.ServicePlanPREMIUMENCRYPTION, // Premium Encryption service plan
	},

	"PurviewDiscovery": {
		constants.SKUSPEE5,                    // Microsoft 365 E5
		constants.SKUPOWERBISTANDARD,          // Power BI Standard (includes Purview Discovery)
		constants.ServicePlanPURVIEWDISCOVERY, // Purview Discovery service plan
	},

	"PurviewMLClassification": {
		constants.SKUSPEE5,                    // Microsoft 365 E5
		constants.ServicePlanMLCLASSIFICATION, // Machine Learning Classification service plan
	},

	"PurviewSafeDocuments": {
		constants.SKUSPEE5,            // Microsoft 365 E5
		constants.ServicePlanSAFEDOCS, // Safe Documents service plan
	},

	// ========================================================================
	// Microsoft 365 Apps & Productivity
	// Reference: https://m365maps.com/ - Microsoft 365 Apps section
	// ========================================================================

	"Microsoft365Apps": {
		constants.SKUSPEE5,                      // Microsoft 365 E5
		constants.ServicePlanOFFICESUBSCRIPTION, // Office Subscription service plan
	},

	"ExchangeOnline": {
		constants.SKUSPEE5,                       // Microsoft 365 E5
		constants.ServicePlanEXCHANGESENTERPRISE, // Exchange Online Plan 2 service plan
		constants.ServicePlanEXCHANGESFOUNDATION, // Exchange Online Foundation service plan
	},

	"ExchangeOnlineArchiving": {
		constants.SKUSPEE5,                     // Microsoft 365 E5
		constants.ServicePlanEXCHANGEANALYTICS, // Exchange Analytics service plan
	},

	"SharePointOnline": {
		constants.SKUSPEE5,                        // Microsoft 365 E5
		constants.ServicePlanSHAREPOINTENTERPRISE, // SharePoint Online Plan 2 service plan
	},

	"OneDriveForBusiness": {
		constants.SKUSPEE5,                        // Microsoft 365 E5
		constants.ServicePlanSHAREPOINTENTERPRISE, // OneDrive included in SharePoint service plan
	},

	"MicrosoftTeams": {
		constants.SKUSPEE5,               // Microsoft 365 E5
		constants.ServicePlanTEAMS1,      // Teams service plan
		constants.ServicePlanMCOSTANDARD, // Teams core features
	},

	"TeamsPhoneSystem": {
		constants.SKUSPEE5,            // Microsoft 365 E5
		constants.ServicePlanMCOEV,    // Teams Phone service plan
		constants.SKUMCOPSTNC,         // Calling Plan SKU
		constants.ServicePlanMCOPSTNC, // Calling Plan service plan
	},

	"TeamsAdvancedMeetings": {
		constants.SKUSPEE5,              // Microsoft 365 E5
		constants.ServicePlanMCOMEETADV, // Advanced Meeting Features service plan
	},

	"TeamsMeshAvatars": {
		constants.SKUSPEE5,                                 // Microsoft 365 E5
		constants.ServicePlanMESHAVATARSFORTEAMS,           // Mesh Avatars for Teams
		constants.ServicePlanMESHAVATARSADDITIONALFORTEAMS, // Mesh Avatars Additional for Teams
	},

	"TeamsMeshImmersive": {
		constants.SKUSPEE5,                         // Microsoft 365 E5
		constants.ServicePlanMESHIMMERSIVEFORTEAMS, // Mesh Immersive for Teams
	},

	"MicrosoftLoop": {
		constants.SKUSPEE5,                 // Microsoft 365 E5
		constants.ServicePlanMICROSOFTLOOP, // Microsoft Loop service plan
	},

	"MicrosoftSearch": {
		constants.SKUSPEE5,                              // Microsoft 365 E5
		constants.ServicePlanMICROSOFTSEARCH,            // Microsoft Search service plan
		constants.ServicePlanGRAPHCONNECTORSSEARCHINDEX, // Graph Connectors Search Index
	},

	"MicrosoftBookings": {
		constants.SKUSPEE5,                     // Microsoft 365 E5
		constants.ServicePlanMICROSOFTBOOKINGS, // Microsoft Bookings service plan
	},

	"MicrosoftForms": {
		constants.SKUSPEE5,               // Microsoft 365 E5
		constants.ServicePlanFORMSPLANE5, // Forms Plan E5 service plan
	},

	"MicrosoftSway": {
		constants.SKUSPEE5,        // Microsoft 365 E5
		constants.ServicePlanSWAY, // Sway service plan
	},

	"MicrosoftClipchamp": {
		constants.SKUSPEE5,             // Microsoft 365 E5
		constants.ServicePlanCLIPCHAMP, // Clipchamp service plan
	},

	"MicrosoftWhiteboard": {
		constants.SKUSPEE5,                   // Microsoft 365 E5
		constants.ServicePlanWHITEBOARDPLAN3, // Whiteboard Plan 3 service plan
	},

	"MicrosoftToDo": {
		constants.SKUSPEE5,              // Microsoft 365 E5
		constants.ServicePlanBPOSSTODO3, // To Do Plan 3 service plan
	},

	"MicrosoftStream": {
		constants.SKUSPEE5,                // Microsoft 365 E5
		constants.ServicePlanSTREAMO365E5, // Stream for Office 365 E5 service plan
	},

	"YammerEnterprise": {
		constants.SKUSPEE5,                    // Microsoft 365 E5
		constants.ServicePlanYAMMERENTERPRISE, // Yammer Enterprise service plan
		constants.ServicePlanVIVAENGAGECORE,   // Viva Engage Core (new Yammer)
	},

	"ProjectForTheWeb": {
		constants.SKUSPEE5,                         // Microsoft 365 E5
		constants.ServicePlanPROJECTO365P3,         // Project for the web P3 service plan
		constants.ServicePlanPROJECTWORKMANAGEMENT, // Project Work Management service plan
	},

	"UniversalPrint": {
		constants.SKUSPEE5,                    // Microsoft 365 E5
		constants.ServicePlanUNIVERSALPRINT01, // Universal Print service plan
	},

	"ExcelPremium": {
		constants.SKUSPEE5,                // Microsoft 365 E5
		constants.ServicePlanEXCELPREMIUM, // Excel Premium (Data Types) service plan
	},

	// ========================================================================
	// Power Platform
	// Reference: https://m365maps.com/ - Power Platform section
	// ========================================================================

	"PowerApps": {
		constants.SKUSPEE5,                   // Microsoft 365 E5
		constants.ServicePlanPOWERAPPSO365P3, // Power Apps for Office 365 P3 service plan
	},

	"PowerAutomate": {
		constants.SKUSPEE5,               // Microsoft 365 E5
		constants.SKUFLOWFREE,            // Power Automate Free SKU
		constants.ServicePlanFLOWO365P3,  // Power Automate for Office 365 P3 service plan
		constants.ServicePlanFLOWP2VIRAL, // Power Automate P2 viral service plan
	},

	"PowerVirtualAgents": {
		constants.SKUSPEE5, // Microsoft 365 E5
		constants.ServicePlanPOWERVIRTUALAGENTSO365P3, // Power Virtual Agents for Office 365 P3
	},

	"Dataverse": {
		constants.SKUSPEE5,                   // Microsoft 365 E5
		constants.ServicePlanCDSO365P3,       // Dataverse for Office 365 P3 service plan
		constants.ServicePlanDYN365CDSO365P3, // Dynamics 365 CDS for Office 365 P3
		constants.ServicePlanDYN365CDSVIRAL,  // Dynamics 365 CDS viral service plan
	},

	// ========================================================================
	// Viva & Analytics
	// Reference: https://m365maps.com/ - Viva section
	// ========================================================================

	"VivaInsights": {
		constants.SKUSPEE5,                            // Microsoft 365 E5
		constants.ServicePlanINSIGHTSBYMYANALYTICS,    // Insights by MyAnalytics service plan
		constants.ServicePlanMICROSOFTMYANALYTICSFULL, // Microsoft MyAnalytics Full service plan
		constants.ServicePlanMYANALYTICSP2,            // MyAnalytics P2 service plan
	},

	"VivaLearning": {
		constants.SKUSPEE5,                      // Microsoft 365 E5
		constants.ServicePlanVIVALEARNINGSEEDED, // Viva Learning Seeded service plan
	},

	"VivaNucleus": {
		constants.SKUSPEE5,           // Microsoft 365 E5
		constants.ServicePlanNucleus, // Nucleus (Viva) service plan
	},

	"VivaEngageCore": {
		constants.SKUSPEE5,                  // Microsoft 365 E5
		constants.ServicePlanVIVAENGAGECORE, // Viva Engage Core service plan
	},

	"VivaPeopleSkills": {
		constants.SKUSPEE5,                          // Microsoft 365 E5
		constants.ServicePlanPEOPLESKILLSFOUNDATION, // People Skills Foundation service plan
	},

	"VivaPlaces": {
		constants.SKUSPEE5,              // Microsoft 365 E5
		constants.ServicePlanPLACESCORE, // Places Core service plan
	},

	// ========================================================================
	// Power BI
	// Reference: https://m365maps.com/ - Power BI section
	// ========================================================================

	"PowerBIStandard": {
		constants.SKUPOWERBISTANDARD,   // Power BI Standard SKU
		constants.ServicePlanBIAZUREP0, // Power BI (free) service plan
	},

	"PowerBIPremium": {
		constants.SKUSPEE5,             // Microsoft 365 E5
		constants.ServicePlanBIAZUREP2, // Power BI Pro service plan
	},

	// ========================================================================
	// Bing Chat Enterprise (Microsoft 365 Copilot Infrastructure)
	// Reference: https://m365maps.com/ - Copilot section
	// ========================================================================

	"BingChatEnterprise": {
		constants.SKUSPEE5,                      // Microsoft 365 E5
		constants.ServicePlanBingChatEnterprise, // Bing Chat Enterprise service plan
	},

	// ========================================================================
	// Microsoft 365 Lighthouse (MSP Tooling)
	// ========================================================================

	"Microsoft365LighthouseCustomer": {
		constants.SKUSPEE5,                               // Microsoft 365 E5
		constants.SKUCPCB2C8RAM128GB,                     // Windows 365 Business
		constants.ServicePlanM365LIGHTHOUSECUSTOMERPLAN1, // M365 Lighthouse Customer Plan 1
	},

	"Microsoft365LighthousePartner": {
		constants.SKUCPCB2C8RAM128GB,                    // Windows 365 Business
		constants.ServicePlanM365LIGHTHOUSEPARTNERPLAN1, // M365 Lighthouse Partner Plan 1
	},
}

// GetRequiredLicensesForFeature returns the list of licenses that enable a specific feature
func GetRequiredLicensesForFeature(featureName string) []string {
	if licenses, ok := RequiredLicenses[featureName]; ok {
		return licenses
	}
	return []string{}
}

// FormatRequiredLicensesMessage returns a formatted error message listing required licenses
func FormatRequiredLicensesMessage(featureName string) string {
	licenses := GetRequiredLicensesForFeature(featureName)
	if len(licenses) == 0 {
		return ""
	}

	message := "This feature requires one of the following licenses:\n"
	for _, lic := range licenses {
		message += "  - " + lic + "\n"
	}
	message += "\nFor more information, see:\n"
	message += "  - https://learn.microsoft.com/en-us/entra/fundamentals/licensing\n"
	message += "  - https://m365maps.com/ (comprehensive licensing diagrams)\n"

	return message
}
