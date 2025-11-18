// Auto-generated from Microsoft Graph API
// Generated: 2025-11-18 12:32:00

package constants

// ============================================================================
// SKU Part Numbers (11 total)
// ============================================================================
const (
	SKUCPCB2C8RAM128GB          = "CPC_B_2C_8RAM_128GB"
	SKUCPCE2C8GB128GB           = "CPC_E_2C_8GB_128GB"
	SKUFLOWFREE                 = "FLOW_FREE"
	SKUMCOPSTNC                 = "MCOPSTNC"
	SKUMicrosoftEntraSuite      = "Microsoft_Entra_Suite"
	SKUMicrosoftIntuneSuite     = "Microsoft_Intune_Suite"
	SKUPOWERBISTANDARD          = "POWER_BI_STANDARD"
	SKURMSBASIC                 = "RMSBASIC"
	SKUSPEE5                    = "SPE_E5"
	SKUWindows365S2vCPU8GB128GB = "Windows_365_S_2vCPU_8GB_128GB"
	SKUWINDOWSSTORE             = "WINDOWS_STORE"
)

// ============================================================================
// Shared Service Plans (5 plans)
// These service plans appear in multiple SKUs
// ============================================================================
const (
	ServicePlanAADPREMIUMP2                = "AAD_PREMIUM_P2"                 // Shared: Microsoft_Entra_Suite, SPE_E5
	ServicePlanEXCHANGESFOUNDATION         = "EXCHANGE_S_FOUNDATION"          // Shared: CPC_B_2C_8RAM_128GB, CPC_E_2C_8GB_128GB, FLOW_FREE, POWER_BI_STANDARD, RMSBASIC, WINDOWS_STORE
	ServicePlanM365LIGHTHOUSECUSTOMERPLAN1 = "M365_LIGHTHOUSE_CUSTOMER_PLAN1" // Shared: CPC_B_2C_8RAM_128GB, SPE_E5
	ServicePlanPURVIEWDISCOVERY            = "PURVIEW_DISCOVERY"              // Shared: POWER_BI_STANDARD, SPE_E5
	ServicePlanWindows10ESUCommercial      = "Windows_10_ESU_Commercial"      // Shared: CPC_B_2C_8RAM_128GB, CPC_E_2C_8GB_128GB
)

// ============================================================================
// Service Plans from: CPC_B_2C_8RAM_128GB (2 plans)
// ============================================================================
const (
	ServicePlanCPCSS2                     = "CPC_SS_2"
	ServicePlanM365LIGHTHOUSEPARTNERPLAN1 = "M365_LIGHTHOUSE_PARTNER_PLAN1"
)

// ============================================================================
// Service Plans from: CPC_E_2C_8GB_128GB (1 plans)
// ============================================================================
const (
	ServicePlanCPC2 = "CPC_2"
)

// ============================================================================
// Service Plans from: FLOW_FREE (2 plans)
// ============================================================================
const (
	ServicePlanDYN365CDSVIRAL = "DYN365_CDS_VIRAL"
	ServicePlanFLOWP2VIRAL    = "FLOW_P2_VIRAL"
)

// ============================================================================
// Service Plans from: MCOPSTNC (1 plans)
// ============================================================================
const (
	ServicePlanMCOPSTNC = "MCOPSTNC"
)

// ============================================================================
// Service Plans from: Microsoft_Entra_Suite (4 plans)
// ============================================================================
const (
	ServicePlanEntraIDentityGovernance             = "Entra_Identity_Governance"
	ServicePlanEntraPremiumInternetAccess          = "Entra_Premium_Internet_Access"
	ServicePlanEntraPremiumPrivateAccess           = "Entra_Premium_Private_Access"
	ServicePlanVerifiableCredentialsServiceRequest = "Verifiable_Credentials_Service_Request"
)

// ============================================================================
// Service Plans from: Microsoft_Intune_Suite (8 plans)
// ============================================================================
const (
	ServicePlan3PARTYAPPPATCH   = "3_PARTY_APP_PATCH"
	ServicePlanCLOUDPKI         = "CLOUD_PKI"
	ServicePlanIntuneAdvancedea = "Intune_AdvancedEA"
	ServicePlanINTUNEP2         = "INTUNE_P2"
	ServicePlanIntuneServicenow = "Intune_ServiceNow"
	ServicePlanIntuneEPM        = "Intune-EPM"
	ServicePlanIntuneMamtunnel  = "Intune-MAMTunnel"
	ServicePlanREMOTEHELP       = "REMOTE_HELP"
)

// ============================================================================
// Service Plans from: POWER_BI_STANDARD (1 plans)
// ============================================================================
const (
	ServicePlanBIAZUREP0 = "BI_AZURE_P0"
)

// ============================================================================
// Service Plans from: RMSBASIC (1 plans)
// ============================================================================
const (
	ServicePlanRMSSBASIC = "RMS_S_BASIC"
)

// ============================================================================
// Service Plans from: SPE_E5 (88 plans)
// ============================================================================
const (
	ServicePlanAADPREMIUM                                = "AAD_PREMIUM"
	ServicePlanADALLOMSO365                              = "ADALLOM_S_O365"
	ServicePlanADALLOMSSTANDALONE                        = "ADALLOM_S_STANDALONE"
	ServicePlanATA                                       = "ATA"
	ServicePlanATPENTERPRISE                             = "ATP_ENTERPRISE"
	ServicePlanBIAZUREP2                                 = "BI_AZURE_P2"
	ServicePlanBingChatEnterprise                        = "Bing_Chat_Enterprise"
	ServicePlanBPOSSTODO3                                = "BPOS_S_TODO_3"
	ServicePlanCDSO365P3                                 = "CDS_O365_P3"
	ServicePlanCLIPCHAMP                                 = "CLIPCHAMP"
	ServicePlanCOMMONDEFENDERPLATFORMFOROFFICE           = "COMMON_DEFENDER_PLATFORM_FOR_OFFICE"
	ServicePlanCOMMUNICATIONSCOMPLIANCE                  = "COMMUNICATIONS_COMPLIANCE"
	ServicePlanCOMMUNICATIONSDLP                         = "COMMUNICATIONS_DLP"
	ServicePlanContentExplorer                           = "Content_Explorer"
	ServicePlanContentexplorerStandard                   = "ContentExplorer_Standard"
	ServicePlanCUSTOMERKEY                               = "CUSTOMER_KEY"
	ServicePlanCustomerlockboxaEnterprise                = "CustomerLockboxA_Enterprise"
	ServicePlanDATAINVESTIGATIONS                        = "DATA_INVESTIGATIONS"
	ServicePlanDefenderForIotEnterprise                  = "Defender_for_Iot_Enterprise"
	ServicePlanDeskless                                  = "Deskless"
	ServicePlanDYN365CDSO365P3                           = "DYN365_CDS_O365_P3"
	ServicePlanEQUIVIOANALYTICS                          = "EQUIVIO_ANALYTICS"
	ServicePlanEXCELPREMIUM                              = "EXCEL_PREMIUM"
	ServicePlanEXCHANGEANALYTICS                         = "EXCHANGE_ANALYTICS"
	ServicePlanEXCHANGESENTERPRISE                       = "EXCHANGE_S_ENTERPRISE"
	ServicePlanFLOWO365P3                                = "FLOW_O365_P3"
	ServicePlanFORMSPLANE5                               = "FORMS_PLAN_E5"
	ServicePlanGRAPHCONNECTORSSEARCHINDEX                = "GRAPH_CONNECTORS_SEARCH_INDEX"
	ServicePlanINFOGOVERNANCE                            = "INFO_GOVERNANCE"
	ServicePlanINFORMATIONBARRIERS                       = "INFORMATION_BARRIERS"
	ServicePlanINSIDERRISK                               = "INSIDER_RISK"
	ServicePlanINSIDERRISKMANAGEMENT                     = "INSIDER_RISK_MANAGEMENT"
	ServicePlanINSIGHTSBYMYANALYTICS                     = "INSIGHTS_BY_MYANALYTICS"
	ServicePlanINTUNEA                                   = "INTUNE_A"
	ServicePlanINTUNEO365                                = "INTUNE_O365"
	ServicePlanKAIZALASTANDALONE                         = "KAIZALA_STANDALONE"
	ServicePlanLOCKBOXENTERPRISE                         = "LOCKBOX_ENTERPRISE"
	ServicePlanM365ADVANCEDAUDITING                      = "M365_ADVANCED_AUDITING"
	ServicePlanM365AUDITPLATFORM                         = "M365_AUDIT_PLATFORM"
	ServicePlanMCOEV                                     = "MCOEV"
	ServicePlanMCOMEETADV                                = "MCOMEETADV"
	ServicePlanMCOSTANDARD                               = "MCOSTANDARD"
	ServicePlanMESHAVATARSADDITIONALFORTEAMS             = "MESH_AVATARS_ADDITIONAL_FOR_TEAMS"
	ServicePlanMESHAVATARSFORTEAMS                       = "MESH_AVATARS_FOR_TEAMS"
	ServicePlanMESHIMMERSIVEFORTEAMS                     = "MESH_IMMERSIVE_FOR_TEAMS"
	ServicePlanMFAPREMIUM                                = "MFA_PREMIUM"
	ServicePlanMICROSOFTCOMMUNICATIONCOMPLIANCE          = "MICROSOFT_COMMUNICATION_COMPLIANCE"
	ServicePlanMICROSOFTLOOP                             = "MICROSOFT_LOOP"
	ServicePlanMICROSOFTMYANALYTICSFULL                  = "MICROSOFT_MYANALYTICS_FULL"
	ServicePlanMICROSOFTSEARCH                           = "MICROSOFT_SEARCH"
	ServicePlanMICROSOFTBOOKINGS                         = "MICROSOFTBOOKINGS"
	ServicePlanMICROSOFTENDPOINTDLP                      = "MICROSOFTENDPOINTDLP"
	ServicePlanMIPSCLP1                                  = "MIP_S_CLP1"
	ServicePlanMIPSCLP2                                  = "MIP_S_CLP2"
	ServicePlanMIPSExchange                              = "MIP_S_Exchange"
	ServicePlanMLCLASSIFICATION                          = "ML_CLASSIFICATION"
	ServicePlanMTP                                       = "MTP"
	ServicePlanMYANALYTICSP2                             = "MYANALYTICS_P2"
	ServicePlanNucleus                                   = "Nucleus"
	ServicePlanOFFICESUBSCRIPTION                        = "OFFICESUBSCRIPTION"
	ServicePlanPAMENTERPRISE                             = "PAM_ENTERPRISE"
	ServicePlanPEOPLESKILLSFOUNDATION                    = "PEOPLE_SKILLS_FOUNDATION"
	ServicePlanPLACESCORE                                = "PLACES_CORE"
	ServicePlanPOWERVIRTUALAGENTSO365P3                  = "POWER_VIRTUAL_AGENTS_O365_P3"
	ServicePlanPOWERAPPSO365P3                           = "POWERAPPS_O365_P3"
	ServicePlanPREMIUMENCRYPTION                         = "PREMIUM_ENCRYPTION"
	ServicePlanPROJECTO365P3                             = "PROJECT_O365_P3"
	ServicePlanPROJECTWORKMANAGEMENT                     = "PROJECTWORKMANAGEMENT"
	ServicePlanRECORDSMANAGEMENT                         = "RECORDS_MANAGEMENT"
	ServicePlanRMSSENTERPRISE                            = "RMS_S_ENTERPRISE"
	ServicePlanRMSSPREMIUM                               = "RMS_S_PREMIUM"
	ServicePlanRMSSPREMIUM2                              = "RMS_S_PREMIUM2"
	ServicePlanSAFEDOCS                                  = "SAFEDOCS"
	ServicePlanSHAREPOINTENTERPRISE                      = "SHAREPOINTENTERPRISE"
	ServicePlanSHAREPOINTWAC                             = "SHAREPOINTWAC"
	ServicePlanSTREAMO365E5                              = "STREAM_O365_E5"
	ServicePlanSWAY                                      = "SWAY"
	ServicePlanTEAMS1                                    = "TEAMS1"
	ServicePlanTHREATINTELLIGENCE                        = "THREAT_INTELLIGENCE"
	ServicePlanUNIVERSALPRINT01                          = "UNIVERSAL_PRINT_01"
	ServicePlanVIVALEARNINGSEEDED                        = "VIVA_LEARNING_SEEDED"
	ServicePlanVIVAENGAGECORE                            = "VIVAENGAGE_CORE"
	ServicePlanWHITEBOARDPLAN3                           = "WHITEBOARD_PLAN3"
	ServicePlanWIN10PROENTSUB                            = "WIN10_PRO_ENT_SUB"
	ServicePlanWINDEFATP                                 = "WINDEFATP"
	ServicePlanWindowsAutopatch                          = "Windows_Autopatch"
	ServicePlanWINDOWSUPDATEFORBUSINESSDEPLOYMENTSERVICE = "WINDOWSUPDATEFORBUSINESS_DEPLOYMENTSERVICE"
	ServicePlanYAMMERENTERPRISE                          = "YAMMER_ENTERPRISE"
)

// ============================================================================
// Service Plans from: Windows_365_S_2vCPU_8GB_128GB (2 plans)
// ============================================================================
const (
	ServicePlanCPCS2C8GB128GB     = "CPC_S_2C_8GB_128GB"
	ServicePlanWINDOWS10ESUTENANT = "WINDOWS_10_ESU_TENANT"
)

// ============================================================================
// Service Plans from: WINDOWS_STORE (1 plans)
// ============================================================================
const (
	ServicePlanWINDOWSSTORE = "WINDOWS_STORE"
)
