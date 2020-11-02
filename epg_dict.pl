#!/usr/bin/env perl
use strict;

use Data::Dumper;
use Convert::ASN1;
use Convert::ASN1 qw(:io :debug);

my @T_PFX = qw( U A C P );

my %ActTime = (
    1  => 'Duration'
   ,2  => 'Inactivity Included'
   ,3  => 'Inactivity'
   ,4  => 'Active Periods'
);

my %APNSel = (
   0 => 'User Equipment or network provided APN, subscription verified'
  ,1 => 'User Equipment provided APN, subscription not verified'
  ,2 => 'Network provided APN, subscription not verified'
);

my %CauseRecClos = (
    0  => 'Normal Release'
   ,4  => 'Abnormal Release'
  ,16  => 'Volume Limit'
  ,17  => 'Time Limit'
  ,18  => 'Serving Node Change'
  ,19  => 'Max Change Condition'
  ,22  => 'RAT Change'
  ,23  => 'TimeZone Change'
  ,24  => 'SGSN PLMN ID Change'
  ,100 => 'Management Init Release'
  ,101 => 'PLMN Change'
  ,102 => 'Credit Control Change'
  ,104 => 'Credit Control Init Release'
  ,105 => 'Policy Control Init Release'
);

my %ChCond = (
   0 => 'QoS Change'
  ,1 => 'Tariff Time'
  ,2 => 'Record Closure'
  ,3 => 'Failure Handling Continue Ongoing'
  ,4 => 'Failure Handling Retry And Terminate Ongoing'
  ,5 => 'Failure Handling Terrminate Ongoing'
);

my %ChSelMode = (
   0   => 'Serving Node Supplied'
  ,3   => 'Home Default'
  ,4   => 'Roaming Default'
  ,5   => 'Visiting Default'
  ,100 => 'Radius Supplied'
  ,101 => 'Roaming Class Based'
);

my %DST = (
  '00' => 'No DST adjustment'
 ,'01' => '+1 hour DST adjustment'
 ,'02' => '+2 hour DST adjustment'
);

my %SrvCondChg = (
   0  => 'QoS Change'
  ,1  => 'SGSN Change'
  ,2  => 'SGSN PLMN Id Change'
  ,3  => 'Tariff Time Switch'
  ,4  => 'PDP Context Release'
  ,5  => 'RAT Change'
  ,6  => 'Service Idle Out'
  ,9  => 'Service Stop'
  ,10 => 'DCCA Time Treshold Reached'
  ,11 => 'DCCA Volume Treshold Reached'
  ,12 => 'DCCA Service Specific Time Treshold Reached'
  ,13 => 'DCCA Time Exhausted'
  ,14 => 'DCCA Volume Exhausted'
  ,15 => 'DCCA Validity Timeout'
  ,17 => 'DCCA Reauthorisation Request'
  ,18 => 'DCCA Continue Ongoing Session'
  ,19 => 'DCCA Retry And Terminate Ongoing Session'
  ,20 => 'DCCA Terminate Ongoing Session'
  ,21 => 'CGI/SAI Change'
  ,22 => 'RAI Change'
  ,23 => 'DCCA Service Specific Unit Exhausted'
  ,24 => 'Record Closure'
  ,29 => 'ECGI Change'
  ,30 => 'TAI Change'

);

my %PCharg = (
  '0100' => 'Hot Billing'
 ,'0200' => 'Flat Rate'
 ,'0400' => 'Prepaid'
 ,'0800' => 'Normal'
);

my %PDNTypeNumber = (
   1 => 'IPv4'
  ,2 => 'IPv6'
  ,4 => 'IPv4v6'
);

my %RATTypes =(
   0 => '<reserved>'
  ,1 => 'UTRAN'
  ,2 => 'GERAN'
  ,3 => 'WLAN'
  ,4 => 'GAN 4'
  ,5 => 'HSPA Evolution'
  ,6 => 'EUTRAN'
);

my %RecTypes = (
  19 => 'GGSN PDP Record'
 ,70 => 'EGSN PDP Record'
 ,84 => 'SGW Record'
 ,85 => 'PGW Record'
);

my %ServingNodeTypes = (
   0 => 'SGSN'
  ,1 => 'PMIPSGW'
  ,2 => 'GTPSGW'
  ,3 => 'EPDG'
  ,4 => 'HSGW'
  ,5 => 'MME'
);

my %YesNoFlag = (
   0 => 'No'
  ,1 => 'Yes'
);

my %TON = (
   0 => '0-unknown'
  ,1 => '1-International'
  ,2 => '2-National'
  ,3 => '3-Network specific'
  ,4 => '4-Subscriber'
  ,5 => '5-reserved'
  ,6 => '6-abbreviated'
  ,7 => '7-reserved for extensionSubscriber'
);

my %NPI = (
   0 => '0-unknown'
  ,1 => '1-ISDN/Telephony'
  ,2 => '2-spare'
  ,3 => '3-Data'
  ,4 => '4-Telex'
  ,5 => '5-spare'
  ,6 => '6-land mobile'
  ,7 => '7-spare'
  ,8 => '8-national'
  ,9 => '9-private'
  ,a => 'A-reserved'
  ,b => 'B-reserved'
  ,c => 'C-reserved'
  ,d => 'D-reserved'
  ,e => 'E-reserved'
  ,f => 'F-reserved for extension'
);

my %Enums = (
   apnSelectionMode          => \%APNSel
  ,causeForRecClosing        => \%CauseRecClos
  ,chargingCharacteristics   => \%PCharg
  ,chChSelectionMode         => \%ChSelMode
  ,changeCondition           => \%ChCond
  ,dynamicAddressFlag        => \%YesNoFlag
  ,rATType                   => \%RATTypes
  ,recordType                => \%RecTypes
  ,servingNodeType           => \%ServingNodeTypes
);

my %Maps = (
   iPBinV4Address            => \&IPdecode
  ,changeTime                => \&Timestmp
  ,eventTimeStamps           => \&Timestmp
  ,mSTimeZone                => \&TimeZone
  ,p_GWPLMNIdentifier        => \&PLMN
  ,pLMNIdentifier            => \&PLMN
  ,pdpPDNType                => \&PDNType
  ,recordOpeningTime         => \&Timestmp
  ,servedIMEISV              => \&TBCD
  ,servedIMSI                => \&TBCD
  ,servedMSISDN              => \&Address
  ,serviceConditionChange    => \&ServCondChg
  ,servingNodePLMNIdentifier => \&PLMN
  ,sgsnPLMNIdentifier        => \&PLMN
  ,sGSNPLMNIdentifier        => \&PLMN
  ,startTime                 => \&Timestmp
  ,stopTime                  => \&Timestmp
  ,timeOfFirstUsage          => \&Timestmp
  ,timeOfLastUsage           => \&Timestmp
  ,timeOfReport              => \&Timestmp
  ,userLocationInformation   => \&UsrLoc
);

my @Zones = (0, 15, 30, 45) ;

my $fh;
my $key;
my $value;
my $buffer;
my $buf1;
my $ofs = 0;

if (@ARGV < 1)
{
  print "Decodes LTE (CDRF_UGW) ASN.1 file to STDOUT\n";
  print "Syntax: ltedecode.pl <inputfile>\n";
  exit;
}

my $inputfile = shift;
my $grammar = << 'EOS' ;
-- EPG DEFINITIONS IMPLICIT TAGS ::=
-- BEGIN
	GPRSRecord ::= CHOICE
	{
		ggsnPDPRecord [21] GGSNPDPRecord,
		egsnPDPRecord [70] EGSNPDPRecord,
		sGWRecord     [78] SGWRecord,
		pgwRecord     [79] PGWRecord
	}
	GGSNPDPRecord ::= SET
	{
		recordType [0] RecordType,
		servedIMSI [3] IMSI,
		ggsnAddress [4] GSNAddress,
		chargingID [5] ChargingID,
		sgsnAddress [6] SEQUENCE OF GSNAddress OPTIONAL,
		accessPointNameNI [7] AccessPointNameNI OPTIONAL,
		pdpType [8] PDPType OPTIONAL,
		servedPDPAddress [9] PDPAddress OPTIONAL,
		dynamicAddressFlag [11] DynamicAddressFlag OPTIONAL,
		listOfTrafficVolumes [12] SEQUENCE OF ChangeOfCharCondition OPTIONAL,
		recordOpeningTime [13] TimeStamp,
		duration [14] CallDuration,
		causeForRecClosing [15] CauseForRecClosing,
		recordSequenceNumber [17] INTEGER OPTIONAL,
		nodeID [18] NodeID OPTIONAL,
		recordExtensions [19] ManagementExtensions OPTIONAL,
		localSequenceNumber [20] LocalSequenceNumber OPTIONAL,
		apnSelectionMode [21] APNSelectionMode OPTIONAL,
		servedMSISDN [22] MSISDN OPTIONAL,
		chargingCharacteristics [23] ChargingCharacteristics,
		chChSelectionMode [24] ChChSelectionMode OPTIONAL,
		iMSsignalingContext [25] NULL OPTIONAL,
		sgsnPLMNIdentifier [27] PLMN-Id OPTIONAL,
		servedIMEISV [29] IMEI OPTIONAL,
		rATType [30] RATType OPTIONAL,
		mSTimeZone [31] MSTimeZone OPTIONAL,
		userLocationInformation [32] BCDString OPTIONAL
--		userLocationInformation [32] OCTET STRING OPTIONAL
	}
	EGSNPDPRecord ::= SET
	{
		recordType [0] RecordType,
		servedIMSI [3] IMSI,
		ggsnAddress [4] GSNAddress,
		chargingID [5] ChargingID,
		sgsnAddress [6] SEQUENCE OF GSNAddress OPTIONAL,
		accessPointNameNI [7] AccessPointNameNI OPTIONAL,
		pdpType [8] PDPType OPTIONAL,
		servedPDPAddress [9] PDPAddress OPTIONAL,
		dynamicAddressFlag [11] DynamicAddressFlag OPTIONAL,
		listOfTrafficVolumes [12] SEQUENCE OF ChangeOfCharCondition OPTIONAL,
		recordOpeningTime [13] TimeStamp,
		duration [14] CallDuration,
		causeForRecClosing [15] CauseForRecClosing,
		recordSequenceNumber [17] INTEGER OPTIONAL,
		nodeID [18] NodeID OPTIONAL,
		recordExtensions [19] ManagementExtensions OPTIONAL,
		localSequenceNumber [20] LocalSequenceNumber OPTIONAL,
		apnSelectionMode [21] APNSelectionMode OPTIONAL,
		servedMSISDN [22] MSISDN OPTIONAL,
		chargingCharacteristics [23] ChargingCharacteristics,
		chChSelectionMode [24] ChChSelectionMode OPTIONAL,
		iMSsignalingContext [25] NULL OPTIONAL,
		sgsnPLMNIdentifier [27] PLMN-Id OPTIONAL,
		pSFurnishChargingInformation [28] PSFurnishChargingInformation OPTIONAL,
		servedIMEISV [29] IMEI OPTIONAL,
		rATType [30] RATType OPTIONAL,
		mSTimeZone [31] MSTimeZone OPTIONAL,
		userLocationInformation [32] BCDString OPTIONAL,
--		userLocationInformation [32] OCTET STRING OPTIONAL,
		listOfServiceData [34] SEQUENCE OF ChangeOfServiceCondition OPTIONAL
	}
	SGWRecord ::= SET
	{
		recordType [0] RecordType,
		servedIMSI [3] IMSI,
		s_GWAddress [4] GSNAddress,
		chargingID [5] ChargingID,
		servingNodeAddress [6] SEQUENCE OF GSNAddress,
		accessPointNameNI [7] AccessPointNameNI OPTIONAL,
		pdpPDNType [8] PDPType OPTIONAL,
		servedPDPPDNAddress [9] PDPAddress OPTIONAL,
		listOfTrafficVolumes [12] SEQUENCE OF ChangeOfCharCondition_SGW OPTIONAL,
		recordOpeningTime [13] TimeStamp,
		duration [14] CallDuration,
		causeForRecClosing [15] CauseForRecClosing,
		recordSequenceNumber [17] INTEGER,
		nodeID [18] NodeID OPTIONAL,
		localSequenceNumber [20] LocalSequenceNumber,
		servedMSISDN [22] MSISDN OPTIONAL,
		chargingCharacteristics [23] ChargingCharacteristics,
		servingNodePLMNIdentifier [27] PLMN-Id OPTIONAL,
		servedIMEISV [29] IMEI OPTIONAL,
		rATType [30] RATType OPTIONAL,
		mSTimeZone [31] MSTimeZone OPTIONAL,
		sGWChange [34] SGWChange OPTIONAL,
		servingNodeType [35] SEQUENCE OF ServingNodeType,
		p_GWAddressUsed [36] GSNAddress OPTIONAL,
		p_GWPLMNIdentifier [37] PLMN-Id OPTIONAL,
		pDNConnectionID [40] ChargingID OPTIONAL,
		servedPDPPDNAddressExt [43] PDPAddress OPTIONAL
	}
	PGWRecord ::= SET
	{
		recordType [0] RecordType,
		servedIMSI [3] IMSI,
		p_GWAddress [4] GSNAddress,
		chargingID [5] ChargingID,
		servingNodeAddress [6] SEQUENCE OF GSNAddress,
		accessPointNameNI [7] AccessPointNameNI OPTIONAL,
		pdpPDNType [8] PDPType OPTIONAL,
		servedPDPPDNAddress [9] PDPAddress OPTIONAL,
		dynamicAddressFlag [11] DynamicAddressFlag OPTIONAL,
		recordOpeningTime [13] TimeStamp,
		duration [14] CallDuration,
		causeForRecClosing [15] CauseForRecClosing,
		recordSequenceNumber [17] INTEGER,
		nodeID [18] NodeID OPTIONAL,
		recordExtensions [19] ManagementExtensions OPTIONAL,
		localSequenceNumber [20] LocalSequenceNumber,
		apnSelectionMode [21] APNSelectionMode OPTIONAL,
		servedMSISDN [22] MSISDN OPTIONAL,
		chargingCharacteristics [23] ChargingCharacteristics,
		chChSelectionMode [24] ChChSelectionMode OPTIONAL,
		iMSsignalingContext [25] NULL OPTIONAL,
		servingNodePLMNIdentifier [27] PLMN-Id OPTIONAL,
		pSFurnishChargingInformation [28] PSFurnishChargingInformation OPTIONAL,
		servedIMEISV [29] IMEI OPTIONAL,
		rATType [30] RATType OPTIONAL,
		mSTimeZone [31] MSTimeZone OPTIONAL,
		userLocationInformation [32] BCDString OPTIONAL,
--		userLocationInformation [32] OCTET STRING OPTIONAL,
		listOfServiceData [34] SEQUENCE OF ChangeOfServiceCondition_PGW OPTIONAL,
		servingNodeType [35] SEQUENCE OF ServingNodeType,
		p_GWPLMNIdentifier [37] PLMN-Id OPTIONAL,
		startTime [38] TimeStamp OPTIONAL,
		stopTime [39] TimeStamp OPTIONAL,
		pDNConnectionID [41] ChargingID OPTIONAL,
		servedPDPPDNAddressExt [45] PDPAddress OPTIONAL
	}
	AccessPointNameNI ::= IA5String --(SIZE(1..63))
	ActiveTimeMethod ::= ENUMERATED
	{
		duration (1),
		inactivityIncluded (2),
		inactivity (3),
		activePeriods (4)
	}
	AddressString ::= BCDString --(SIZE (1..20))
--	AddressString ::= OCTET STRING --(SIZE (1..20))
	AFRecordInformation ::= SEQUENCE
	{
		aFChargingIdentifier [1] AFChargingIdentifier
	}
	AFChargingIdentifier ::= OCTET STRING
	APNSelectionMode ::= ENUMERATED
	{
		mSorNetworkProvidedSubscriptionVerified (0),
		mSProvidedSubscriptionNotVerified (1),
		networkProvidedSubscriptionNotVerified (2)
	}
	CallDuration ::= INTEGER
	CauseForRecClosing ::= ENUMERATED
	{
		normalRelease (0),
		abnormalRelease (4),
		volumeLimit (16),
		timeLimit (17),
		sGSNChange (18),
		maxChangeCond (19),
		rATChange (22),
		mSTimeZoneChange (23),
		sGSNPLMNIDChange (24),
		managementInitRelease (100),
		creditControlChange (102),
		creditControlInitRelease (104),
		policyControlInitRelease (105)
	}
	ChangeCondition ::= ENUMERATED
	{
		qoSChange (0),
		tariffTime (1),
		recordClosure (2),
		failureHandlingContinueOngoing (3),
		failureHandlingRetryandTerminateOngoing (4),
		failureHandlingTerminateOngoing (5)
	}
	ChangeOfCharCondition ::= SEQUENCE
	{
		qosNegotiated [2] QoSInformation OPTIONAL,
		dataVolumeGPRSUplink [3] DataVolumeGPRS,
		dataVolumeGPRSDownlink [4] DataVolumeGPRS,
		changeCondition [5] ChangeCondition,
		changeTime [6] TimeStamp,
		userLocationInformation [8] BCDString OPTIONAL
--		userLocationInformation [8] OCTET STRING OPTIONAL
	}
	ChangeOfCharCondition_SGW ::= SEQUENCE
	{
		dataVolumeGPRSUplink [3] DataVolumeGPRS,
		dataVolumeGPRSDownlink [4] DataVolumeGPRS,
		changeCondition [5] ChangeCondition,
		changeTime [6] TimeStamp,
		userLocationInformation [8] BCDString OPTIONAL,
--		userLocationInformation [8] OCTET STRING OPTIONAL,
		ePCQoSInformation [9] EPCQoSInformation OPTIONAL
	}
	ChangeOfServiceCondition ::= SEQUENCE
	{
		ratingGroup [1] RatingGroupId,
		resultCode [3] ResultCode OPTIONAL,
		localSequenceNumber [4] LocalSequenceNumber OPTIONAL,
		timeOfFirstUsage [5] TimeStamp OPTIONAL,
		timeOfLastUsage [6] TimeStamp OPTIONAL,
		timeUsage [7] CallDuration OPTIONAL,
		serviceConditionChange [8] ServiceConditionChange,
		qoSInformationNeg [9] QoSInformation OPTIONAL,
		sgsn-Address [10] GSNAddress OPTIONAL,
		sGSNPLMNIdentifier [11] PLMN-Id OPTIONAL,
		datavolumeFBCUplink [12] DataVolumeGPRS OPTIONAL,
		datavolumeFBCDownlink [13] DataVolumeGPRS OPTIONAL,
		timeOfReport [14] TimeStamp,
		rATType [15] RATType OPTIONAL,
		failureHandlingContinue [16] FailureHandlingContinue OPTIONAL,
		serviceIdentifier [17] ServiceIdentifier OPTIONAL,
		pSFurnishChargingInformation [18] PSFurnishChargingInformation OPTIONAL,
		aFRecordInformation [19] SEQUENCE OF OCTET STRING OPTIONAL,
		userLocationInformation [20] BCDString OPTIONAL,
--		userLocationInformation [20] OCTET STRING OPTIONAL,
		eventBasedChargingInformation [21] EventBasedChargingInformation OPTIONAL
	}
	ChangeOfServiceCondition_PGW ::= SEQUENCE
	{
		ratingGroup [1] RatingGroupId,
		resultCode [3] ResultCode OPTIONAL,
		localSequenceNumber [4] LocalSequenceNumber OPTIONAL,
		timeOfFirstUsage [5] TimeStamp OPTIONAL,
		timeOfLastUsage [6] TimeStamp OPTIONAL,
		timeUsage [7] CallDuration OPTIONAL,
		serviceConditionChange [8] ServiceConditionChange,
		qoSInformationNeg [9] EPCQoSInformation OPTIONAL,
		sgsn-Address [10] GSNAddress OPTIONAL,
		sGSNPLMNIdentifier [11] PLMN-Id OPTIONAL,
		datavolumeFBCUplink [12] DataVolumeGPRS OPTIONAL,
		datavolumeFBCDownlink [13] DataVolumeGPRS OPTIONAL,
		timeOfReport [14] TimeStamp,
		rATType [15] RATType OPTIONAL,
		failureHandlingContinue [16] FailureHandlingContinue OPTIONAL,
		serviceIdentifier [17] ServiceIdentifier OPTIONAL,
		pSFurnishChargingInformation [18] PSFurnishChargingInformation OPTIONAL,
		aFRecordInformation [19] SEQUENCE OF AFRecordInformation OPTIONAL,
		userLocationInformation [20] BCDString OPTIONAL,
--		userLocationInformation [20] OCTET STRING OPTIONAL,
		eventBasedChargingInformation [21] EventBasedChargingInformation OPTIONAL
	}

	ChargingCharacteristics ::= BCDString --(SIZE (2))
--	ChargingCharacteristics ::= OCTET STRING --(SIZE (2))
	ChargingID ::= INTEGER --(0..4294967295)
	ChChSelectionMode ::= ENUMERATED
	{
		sGSNSupplied (0),
		homeDefault (3),
		roamingDefault (4),
		visitingDefault (5),
		radiusSupplied (100),
		roamingClassBased (101)
	}
	CreditControlFailureReport ::= SEQUENCE
	{
		requestType [0] CreditRequestType,
		requestStatus [1] CreditRequestStatus,
		resultCode [2] CreditResultCode OPTIONAL,
		ccRequestNumber [12] INTEGER OPTIONAL
	}
	CreditControlInfo ::= SEQUENCE
	{
		creditControlFailureReport [6] CreditControlFailureReport OPTIONAL,
		creditControlSessionId [7] OCTET STRING OPTIONAL, --(SIZE(1..255)) OPTIONAL,
		ccsRealm [8] OCTET STRING OPTIONAL --(SIZE(1..255)) OPTIONAL
	}
	CreditRequestType ::= ENUMERATED
	{
		start (0),
		interim (1),
		stop (2)
	}
	CreditRequestStatus ::= ENUMERATED
	{
		unsent (0),
		noAnswer (1),
		failure (2)
	}
	CreditResultCode ::= INTEGER
	DataVolumeGPRS ::= INTEGER
	DynamicAddressFlag ::= BOOLEAN
	EPCQoSInformation ::= SEQUENCE
	{
		--
		-- See TS 29.212 for more information
		--
		qCI [1] INTEGER,
		maxRequestedBandwithUL [2] INTEGER OPTIONAL,
		maxRequestedBandwithDL [3] INTEGER OPTIONAL,
		guaranteedBitrateUL [4] INTEGER OPTIONAL,
		guaranteedBitrateDL [5] INTEGER OPTIONAL,
		aRP [6] INTEGER OPTIONAL
	}
	ETSIAddress ::= AddressString
	EventBasedChargingInformation ::= SEQUENCE
	{
		numberOfEvents [1] INTEGER,
		eventTimeStamps [2] SEQUENCE OF TimeStamp OPTIONAL
	}
	FailureHandlingContinue ::= BOOLEAN
	FFDAppendIndicator ::= BOOLEAN
	FreeFormatData ::= OCTET STRING --(SIZE (1..160))
	GprsCdrExtensions ::= SET
	{
		creditControlInfo [2] CreditControlInfo OPTIONAL,
		policyControlInfo [3] PolicyControlInfo OPTIONAL,
		userCategory [5] INTEGER OPTIONAL,
		ruleSpaceId [6] IA5String OPTIONAL,
		serviceContainers [7] SEQUENCE OF ServiceContainer OPTIONAL,
		timeReports [8] SEQUENCE OF TimeReport OPTIONAL
	}
	GSNAddress ::= IPAddress
	IMEI ::= TBCD-STRING --(SIZE (8))
	IMSI ::= TBCD-STRING --(SIZE (3..8))
	IPAddress ::= CHOICE
	{
		iPBinaryAddress IPBinaryAddress,
		iPTextRepresentedAddress IPTextRepresentedAddress
	}
	IPBinaryAddress ::= CHOICE
	{
		iPBinV4Address [0] BCDString, --(SIZE(4)),
--		iPBinV4Address [0] OCTET STRING, --(SIZE(4)),
		iPBinV6Address [1] BCDString --(SIZE(16))
--		iPBinV6Address [1] OCTET STRING --(SIZE(16))
	}
	IPTextRepresentedAddress ::= CHOICE
	{
		iPTextV4Address [2] IA5String, --(SIZE(7..15)),
		iPTextV6Address [3] IA5String --(SIZE(15..45))
	}
	ISDN-AddressString ::= AddressString --(SIZE(1..9))
	LocalSequenceNumber ::= INTEGER --(0..4294967295)
	ManagementExtensions ::= SET OF ManagementExtension
	ManagementExtension ::= SEQUENCE
	{
		identifier OBJECT IDENTIFIER,
		significance [1] BOOLEAN OPTIONAL, --DEFAULT TRUE,
		information [2] GprsCdrExtensions
	}
	MSISDN ::= ISDN-AddressString
	MSTimeZone ::= BCDString --(SIZE (2))
--	MSTimeZone ::= OCTET STRING --(SIZE (2))
	NodeID ::= IA5String --(SIZE(1..20))
	PDPAddress ::= CHOICE
	{
		iPAddress [0] IPAddress,
		eTSIAddress [1] ETSIAddress
	}
	PDPType ::= BCDString --(SIZE(2))
--	PDPType ::= OCTET STRING --(SIZE(2))
	PLMN-Id ::= BCDString --(SIZE(3))
--	PLMN-Id ::= OCTET STRING --(SIZE(3))
	PolicyControlFailureReport ::= SEQUENCE
	{
		requestType [0] PolicyRequestType,
		requestStatus [1] PolicyRequestStatus,
		resultCode [2] PolicyResultCode OPTIONAL,
		stopTime [5] TimeStamp OPTIONAL
	}
	PolicyControlInfo ::= SEQUENCE
	{
		policyControlFailureReport [4] PolicyControlFailureReport OPTIONAL,
		pcsRealm [6] OCTET STRING OPTIONAL, --(SIZE(1..255)) OPTIONAL,
		policyControlSessionId [7] OCTET STRING OPTIONAL --(SIZE(1..255)) OPTIONAL
	}
	PolicyRequestType ::= ENUMERATED
	{
		start (0),
		interim (1),
		stop (2)
	}
	PolicyRequestStatus ::= ENUMERATED
	{
		unsent (0),
		noAnswer (1),
		failure (2)
	}
	PolicyResultCode ::= INTEGER
	PSFurnishChargingInformation ::= SEQUENCE
	{
		pSFreeFormatData [1] FreeFormatData,
		pSFFDAppendIndicator [2] FFDAppendIndicator OPTIONAL
	}
	QoSInformation ::= BCDString --(SIZE (4..15))
--	QoSInformation ::= OCTET STRING --(SIZE (4..15))
	RatingGroupId ::= INTEGER
	RATType ::= INTEGER --(0..255)
	RecordType ::= ENUMERATED
	{
		ggsnPDPRecord (19),
		egsnPDPRecord (70),
		sGWRecord (84),
		pGWRecord (85)
	}
	ResultCode ::= INTEGER
	ServiceConditionChange ::= BCDString
	ServiceContainer ::= SEQUENCE
	{
		ratingGroup [1] RatingGroupId,
		serviceIdentifier [2] ServiceIdentifier OPTIONAL,
		localSequenceNumber [3] LocalSequenceNumber OPTIONAL,
		method [4] ActiveTimeMethod OPTIONAL,
		inactivity [5] INTEGER OPTIONAL,
		resolution [6] INTEGER OPTIONAL,
		ccRequestNumber [7] INTEGER OPTIONAL,
		serviceSpecificUnits [8] INTEGER OPTIONAL,
		listOfURI [9] SEQUENCE OF URI OPTIONAL
	}
	ServiceIdentifier ::= INTEGER --(0..4294967295)
	ServingNodeType ::= ENUMERATED
	{
		sGSN (0),
		pMIPSGW (1),
		gTPSGW (2),
		ePDG (3),
		hSGW (4),
		mME (5)
	}
	SGWChange ::= BOOLEAN
	TBCD-STRING ::= BCDString
--	TBCD-STRING ::= OCTET STRING
	TimeReport ::= SEQUENCE
	{
		ratingGroup [1] RatingGroupId,
		startTime [2] TimeStamp,
		endTime [3] TimeStamp,
		dataVolumeUplink [4] DataVolumeGPRS OPTIONAL,
		dataVolumeDownlink [5] DataVolumeGPRS OPTIONAL
	}
	TimeStamp ::= BCDString --(SIZE(9))
--	TimeStamp ::= OCTET STRING --(SIZE(9))
	URI ::= SEQUENCE
	{
		count [1] INTEGER OPTIONAL,
		uri [2] IA5String OPTIONAL,
		uriIdentifier [3] INTEGER OPTIONAL,
		uriDataVolumeUplink [4] INTEGER OPTIONAL,
		uriDataVolumeDownlink [5] INTEGER OPTIONAL,
		listOfUriTimeStamps [6] SET OF TimeStamp OPTIONAL
	}
-- END
EOS

my $root = 'GPRSRecord';

my $asn;
my $out;
my $bsn = Convert::ASN1->new;
my $retcode = $bsn->prepare($grammar) ;

my $ast = $bsn->find($root);

dump_ast($ast);

die "Done.";


die "Prepare error:\n".$bsn->error."\n"
   unless (defined $retcode) ;
open ($fh, $inputfile);

while (1)
{
  $retcode = asn_read($fh, $buffer);
  last unless defined $retcode ;

  if ($retcode == -1)
  {
     print "Incomplete...\n";
     last ;
  }

  $asn = $bsn->find($root);
  die "Could not find $root!\n"
     unless defined $asn;

  $out = $asn->decode($buffer)
     or die 'error: ' . $asn->error . "\n";

  &show_it($root, $out);
}

close $fh;

die "Read error: \n".$bsn->error."\n"
  if $bsn->error ;

sub show_it {
  my $key = shift ;
  my $value = shift ;
  my $vofs = shift ;

  if (defined $vofs)
  {
    $ofs += 3 ;
  }

  my $filler = ' ' x $ofs;
  my $t = ref $value;

  REF_TYPE:
  {
      $t =~ /ARRAY/                   && do {
                     print "$filler$key:$t\n";
                     for (my $i = 0; $i < @$value; $i++)
                     {
                        my $v = $$value[$i];
                        my $tt = ref $v;
                        $v = &Mapka($key, $v)
                           unless $tt eq 'HASH' or $tt eq 'ARRAY' ;
                        &show_it("[$i]", $v, 1);
                     }
                     last REF_TYPE;
                                              };

      $t =~ /HASH/                   && do {
                     $t .= " ($retcode bytes) "
                        if $key eq 'GPRSCallEventRecord' ;
                     print "$filler$key:$t\n";
                     foreach my $k (sort keys %$value)
                     {
                        my $v = $value->{$k} ;
                        &show_it($k, $v, 1);
                     }
                     last REF_TYPE;
                                              };
      $t =~ /Math::BigInt/           && do {
                     $value =~ s/^\+(\d*)/$1/ ;
                                              };

      1 == 1                         && do {
                     $value = &Mapka($key, $value) ;
                     print "$filler$key=$value\n";
                                              };
  } # REF_TYPE

  $ofs -= 3
    if ($ofs);
}

sub Mapka {
  my $k = shift;
  my $v = shift;

  if (exists $Maps{$k})
  {
    return $Maps{$k}($v, $k) ;
  }
  else
  {
    return $v . &Enum($k, $v) ;
  }
}

sub Enum {
  my $k = shift;
  my $v = shift;

  return ''
    unless (exists $Enums{$k}) ;

  if (exists $Enums{$k}->{$v})
  {
     return ' (' . $Enums{$k}->{$v} . ')' ;
  }
  else
  {
     return ' (Unknown)' ;
  }
}

sub TBCD {
  my $val = shift ;

  return &RevBCD($val, 1) ;
}

sub RevBCD {
  my $val = shift ;
  my $cond = shift ;
  my @arr = () ;

  for (my $i = 0; $i < length($val); $i += 2)
  {
    push @arr, substr($val, $i+1, 1) ;
    push @arr, substr($val, $i, 1) ;
  }
  pop @arr
      if defined $cond && 'f' eq $arr[$#arr] ;
  return join('', @arr) ;
}

sub IPdecode {
  my $val = shift ;
  
  if (length($val) % 2)
  {
    $val = $val . 'f';
  }

  return join('.', map (hex, ($val =~  /(..)/g))) ;
}

sub PLMN {
  my $val = shift ;

  if ('f' eq substr($val, 2, 1))
  {
     return &TBCD(substr($val, 0, 4)) . "-" . &RevBCD(substr($val, 4), 0);
  }
  else
  {
     return &RevBCD(substr($val, 0, 2), 0) . substr($val, 3, 1) . &RevBCD(substr($val, 4), 0) . substr($val, 2, 1);
  }
}

sub Timestmp {
  my $val = shift ;

  my $utc = chr(hex(substr($val, 12, 2))) . substr($val, 14);

  return sprintf("20%02d-%02d-%02d %02d:%02d:%02d [%s]", unpack ('A2A2A2A2A2A2', $val), $utc);
}

sub TimeZone {
  my $val = shift ;
  my $aux = &RevBCD(substr($val, 0, 2), 0);
  my $sign = $aux & 0x80 ;
  my $tzone = $aux & 0x7F ;

  return sprintf("%s%02d%02d", ($sign ? '-' : '+'), ($tzone/4), ($Zones[$tzone%4]))
         . ' [' . &GetMap(\%DST, substr($val, 2, 2))
         . ']' ;
}


sub UsrLoc {
  my $val = shift ;
  my $flags = hex(substr($val, 0, 2)) ;
  my $MASK = 1 ;

  my $ofs = 2 ;
  my $filler = '' ;
  my $outstring = '' ;

  if (0==$flags)
  {
    # CGI only
    $outstring .= $filler . 'CGI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/LAC:0x' . substr($val, $ofs, 4) . '/CI:0x' . substr($val, $ofs+4, 4) ;
    return $outstring ;
  }

  if (1==$flags)
  {
    # SAI only
    $outstring .= $filler . 'SAI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/LAC:0x' . substr($val, $ofs, 4) . '/SAC:0x' . substr($val, $ofs+4, 4) ;
    return $outstring ;
  }

  $flags &= hex('3F') ;

  if ($flags & $MASK)
  {
    # CGI
    $outstring .= $filler . 'CGI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/LAC:0x' . substr($val, $ofs, 4) . '/CI:0x' . substr($val, $ofs+4, 4) ;
    $ofs += 8 ;
  }

  $flags >>= 1 ;
  if ($flags & $MASK)
  {
    # SAI
    $outstring .= $filler . 'SAI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/LAC:0x' . substr($val, $ofs, 4) . '/SAC:0x' . substr($val, $ofs+4, 4) ;
    $ofs += 8 ;
  }

  $flags >>= 1 ;
  if ($flags & $MASK)
  {
    # RAI
    $outstring .= $filler . 'RAI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/LAC:0x' . substr($val, $ofs, 4) . '/RAC:0x' . substr($val, $ofs+4, 2) ;
    $ofs += 8 ;
  }

  $flags >>= 1 ;
  if ($flags & $MASK)
  {
    # TAI
    $outstring .= $filler . 'TAI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/TAC:0x' . substr($val, $ofs, 4) ;
    $ofs += 4 ;
  }

  $flags >>= 1 ;
  if ($flags & $MASK)
  {
    # ECGI
    $outstring .= $filler . 'ECGI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/ECI:0x' . substr($val, $ofs+1, 7) ;
    $ofs += 8 ;
  }

  $flags >>= 1 ;
  if ($flags & $MASK)
  {
    # LAI
    $outstring .= $filler . 'LAI/PLMN:' . &PLMN(substr($val, $ofs, 6)) ;
    $ofs += 6;
    $filler = ' ';

    $outstring .= '/LAC:0x' . substr($val, $ofs, 4) ;
    $ofs += 4 ;
  }

  return $outstring ;
}

sub Address {
  my $val = shift ;

  return RevBCD(substr($val, 2), 0)
         . ' [TON: ' . GetMap(\%TON, 0x07 & substr($val, 0, 1))
         . ', NPI: ' . GetMap(\%NPI, substr($val, 1, 1))
         . ']';
}

sub Bytes {
  my $val = shift ;
  my @arr = ();

  #return '{0x' . join(' 0x', unpack('H2' x int(length ($val) / 2), $val)) . '}' ;

  for (my $i = 0; $i < length($val); $i += 2)
  {
    push @arr, substr($val, $i, 2) ;
  }
  return '{0x' . join(' 0x', @arr) . '}' ;
}

sub PDNType {
  my $val = shift ;
  my $aux = hex (substr($val, 2, 2)) & 0x07;

  return &Bytes($val)
         . ' (' . &GetMap(\%PDNTypeNumber, $aux)
         . ')' ;
}

sub ServCondChg {
  my $val = shift;
  my $retval = '0x' . $val ;
  my @arr = ();
  my $mask = 0x80000000;

  $val = hex($retval);

  for (my $i = 0; $i < 32; $i++, $mask >>= 1)
  {

    if ($val & $mask)
    {
       if (exists $SrvCondChg{$i}) 
       {
         push @arr, $SrvCondChg{$i} ;
       }
       else
       {
         push @arr, "Unknown $i!" ;
       }
    }

  }

  $retval .= ' (' . join (', ', @arr) . ')'
       if (@arr > 0) ;

  return $retval ;
}

sub GetMap {
  my $map = shift;
  my $key = shift;

  if (exists $map->{$key})
  {
     return $map->{$key} ;
  }
  else
  {
     return "Unknown value $key" ;
  }
}

sub dump_ast {
	my $root = shift->{tree}->{GPRSRecord};

	my $ops = {
		0  => 'valUNKNOWN',
		1  => 'valBOOLEAN',
		2  => 'valINTEGER',
		3  => 'valBITSTR',
		4  => 'valSTRING',
		5  => 'valNULL',
		6  => 'valOBJID',
		7  => 'valREAL',
		8  => 'valSEQUENCE',
		9  => 'valEXPLICIT',
		10 => 'valSET',
		11 => 'valUTIME',
		12 => 'valGTIME',
		13 => 'valUTF8',
		14 => 'valANY',
		15 => 'valCHOICE',
		16 => 'valROID',
		17 => 'valBCD',
		18 => 'valEXTENSIONS',
	};

	sub scan {
		my $node = shift;
		my $pfx = shift;

		if (ref $node eq 'ARRAY') {
			for my $pr (@$node) {
				scan_node($pr, $pfx);
			}
		}
		else {
			scan_node($node, $pfx);
		}
	}

	sub scan_node {
		my $pr = shift;
		my $pfx = shift;

		my ($TAG, $TYPE, $VAR, $LOOP, $OPT, $EXT, $CHILD, $DEFINE) = @$pr;
		my $ts = ($TAG ne '') ? tag_str($TAG) : '';

		print qq{dictionary("$pfx.$ts", "$VAR", $ops->{$TYPE});\n} if ($ts && $VAR);
		print qq{dictionary("$pfx", "$VAR", $ops->{$TYPE});\n} if (!$ts && $VAR);

		if ($CHILD) {
			scan($CHILD, ($ts ? "$pfx.$ts" : $pfx));
		}
	}

	sub tag_str {
		my @cc = map { ord } split(//, $_[0]);

		my $c = shift @cc;
		my $res = $T_PFX[($c & 0xFF) >> 6];
		$c &= 0x1F;

		if ($c == 31) {
			$res .= (shift(@cc) & 0x7F);
		} else {
			$res .= $c;
		}
		return $res;
	}

	scan($root, '');
}

