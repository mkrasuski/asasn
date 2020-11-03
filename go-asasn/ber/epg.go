package ber

import "strconv"

// FieldFormater - coverter from Buffer to output for the field
type FieldFormater func(*Buffer, Output, int)

// Field - metadata
type Field struct {
	name   string
	format FieldFormater
}

var epgFields = make(map[string]*Field)
var schemaLoaded = epgPatches(epgSchema())

// FieldByPath -
func FieldByPath(path string) (*Field, bool) {
	f, ok := epgFields[path]
	return f, ok
}

func empty(b *Buffer, out Output, len int) {
	// EMPTY
}

func dictionary(path string, name string, format FieldFormater) {

	epgFields[path] = &Field{
		name:   name,
		format: format,
	}
}

var valSTRING = (*Buffer).AsString
var valNULL = (*Buffer).AsString
var valINTEGER = (*Buffer).AsInt
var valBCD = (*Buffer).AsBCD
var valBOOLEAN = (*Buffer).AsBool
var valCHOICE = empty
var valSET = empty
var valSEQUENCE = empty
var valOBJID = (*Buffer).AsHex

func patchField(name string, fmter FieldFormater) {

	for _, v := range epgFields {
		if v.name == name {
			v.format = fmter
		}
	}
}

func patchManyFields(fmter FieldFormater, fields ...string) {
	for _, v := range epgFields {
		for _, fn := range fields {
			if v.name == fn {
				v.format = fmter
			}
		}
	}
}

func makeEnum(m map[int64]string) FieldFormater {

	return func(b *Buffer, out Output, len int) {
		val := b.readInt(len)
		out.WriteString(strconv.FormatInt(val, 10))
		if label, ok := m[val]; ok {
			out.WriteByte(' ')
			out.WriteString(label)
		}
	}
}

func epgPatches(deps bool) bool {

	patchField("serviceConditionChange", makeEnum(map[int64]string{
		0:  "// QoS Change",
		1:  "// SGSN Change",
		2:  "// SGSN PLMN Id Change",
		3:  "// Tariff Time Switch",
		4:  "// PDP Context Release",
		5:  "// RAT Change",
		6:  "// Service Idle Out",
		9:  "// Service Stop",
		10: "// DCCA Time Treshold Reached",
		11: "// DCCA Volume Treshold Reached",
		12: "// DCCA Service Specific Time Treshold Reached",
		13: "// DCCA Time Exhausted",
		14: "// DCCA Volume Exhausted",
		15: "// DCCA Validity Timeout",
		17: "// DCCA Reauthorisation Request",
		18: "// DCCA Continue Ongoing Session",
		19: "// DCCA Retry And Terminate Ongoing Session",
		20: "// DCCA Terminate Ongoing Session",
		21: "// CGI/SAI Change",
		22: "// RAI Change",
		23: "// DCCA Service Specific Unit Exhausted",
		24: "// Record Closure",
		29: "// ECGI Change",
		30: "// TAI Change",
	}))

	patchField("servingNodeType", makeEnum(map[int64]string{
		0: "// SGSN",
		1: "// PMIPSGW",
		2: "// GTPSGW",
		3: "// EPDG",
		4: "// HSGW",
		5: "// MME",
	}))

	patchField("recordType", makeEnum(map[int64]string{
		19: "// GGSN PDP Record",
		70: "// EGSN PDP Record",
		84: "// SGW Record",
		85: "// PGW Record",
	}))

	patchField("rATType", makeEnum(map[int64]string{
		0: "// <reserved>",
		1: "// UTRAN",
		2: "// GERAN",
		3: "// WLAN",
		4: "// GAN 4",
		5: "// HSPA Evolution",
		6: "// EUTRAN",
	}))

	patchField("causeForRecClosing", makeEnum(map[int64]string{
		0:   "// Normal Release",
		4:   "// Abnormal Release",
		16:  "// Volume Limit",
		17:  "// Time Limit",
		18:  "// Serving Node Change",
		19:  "// Max Change Condition",
		22:  "// RAT Change",
		23:  "// TimeZone Change",
		24:  "// SGSN PLMN ID Change",
		100: "// Management Init Release",
		101: "// PLMN Change",
		102: "// Credit Control Change",
		104: "// Credit Control Init Release",
		105: "// Policy Control Init Release",
	}))

	patchField("apnSelectionMode", makeEnum(map[int64]string{
		0: "// User Equipment or network provided APN, subscription verified",
		1: "// User Equipment provided APN, subscription not verified",
		2: "// Network provided APN, subscription not verified",
	}))

	patchField("chChSelectionMode", makeEnum(map[int64]string{
		0:   "// Serving Node Supplied",
		3:   "// Home Default",
		4:   "// Roaming Default",
		5:   "// Visiting Default",
		100: "// Radius Supplied",
		101: "// Roaming Class Based",
	}))

	patchField("changeCondition", makeEnum(map[int64]string{
		0: "// QoS Change",
		1: "// Tariff Time",
		2: "// Record Closure",
		3: "// Failure Handling Continue Ongoing",
		4: "// Failure Handling Retry And Terminate Ongoing",
		5: "// Failure Handling Terrminate Ongoing",
	}))

	patchField("chargingCharacteristics", makeEnum(map[int64]string{
		0x0100: "// Hot Billing",
		0x0200: "// Flat Rate",
		0x0400: "// Prepaid",
		0x0800: "// Normal",
	}))

	patchManyFields((*Buffer).AsRHex,
		"servedIMEISV",
		"servedIMSI",
	)

	patchManyFields((*Buffer).AsIPAddress,
		"iPBinV4Address",
	)

	patchManyFields((*Buffer).AsTimestamp,
		"changeTime",
		"eventTimeStamps",
		"recordOpeningTime",
		"startTime",
		"stopTime",
		"timeOfFirstUsage",
		"timeOfLastUsage",
		"timeOfReport",
	)

	patchManyFields((*Buffer).AsPLMN,
		"p_GWPLMNIdentifier",
		"pLMNIdentifier",
		"servingNodePLMNIdentifier",
		"sgsnPLMNIdentifier",
	)

	return true
}

func epgSchema() bool {
	dictionary(".C21", "ggsnPDPRecord", valSET)
	dictionary(".C21.C0", "recordType", valINTEGER)
	dictionary(".C21.C3", "servedIMSI", valBCD)
	dictionary(".C21.C5", "chargingID", valINTEGER)
	dictionary(".C21.C7", "accessPointNameNI", valSTRING)
	dictionary(".C21.C8", "pdpType", valBCD)
	dictionary(".C21.C11", "dynamicAddressFlag", valBOOLEAN)
	dictionary(".C21.C13", "recordOpeningTime", valBCD)
	dictionary(".C21.C14", "duration", valINTEGER)
	dictionary(".C21.C15", "causeForRecClosing", valINTEGER)
	dictionary(".C21.C17", "recordSequenceNumber", valINTEGER)
	dictionary(".C21.C18", "nodeID", valSTRING)
	dictionary(".C21.C20", "localSequenceNumber", valINTEGER)
	dictionary(".C21.C21", "apnSelectionMode", valINTEGER)
	dictionary(".C21.C22", "servedMSISDN", valBCD)
	dictionary(".C21.C23", "chargingCharacteristics", valBCD)
	dictionary(".C21.C24", "chChSelectionMode", valINTEGER)
	dictionary(".C21.C25", "iMSsignalingContext", valNULL)
	dictionary(".C21.C27", "sgsnPLMNIdentifier", valBCD)
	dictionary(".C21.C29", "servedIMEISV", valBCD)
	dictionary(".C21.C30", "rATType", valINTEGER)
	dictionary(".C21.C31", "mSTimeZone", valBCD)
	dictionary(".C21.C32", "userLocationInformation", valBCD)
	dictionary(".C21.C4", "ggsnAddress", valSEQUENCE)
	dictionary(".C21.C4", "iPBinaryAddress", valCHOICE)
	dictionary(".C21.C4.C0", "iPBinV4Address", valBCD)
	dictionary(".C21.C4.C1", "iPBinV6Address", valBCD)
	dictionary(".C21.C4", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C21.C4.C2", "iPTextV4Address", valSTRING)
	dictionary(".C21.C4.C3", "iPTextV6Address", valSTRING)
	dictionary(".C21.C6", "sgsnAddress", valSEQUENCE)
	dictionary(".C21.C6", "iPBinaryAddress", valCHOICE)
	dictionary(".C21.C6.C0", "iPBinV4Address", valBCD)
	dictionary(".C21.C6.C1", "iPBinV6Address", valBCD)
	dictionary(".C21.C6", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C21.C6.C2", "iPTextV4Address", valSTRING)
	dictionary(".C21.C6.C3", "iPTextV6Address", valSTRING)
	dictionary(".C21.C9", "servedPDPAddress", valSEQUENCE)
	dictionary(".C21.C9.C0", "iPAddress", valSEQUENCE)
	dictionary(".C21.C9.C0", "iPBinaryAddress", valCHOICE)
	dictionary(".C21.C9.C0.C0", "iPBinV4Address", valBCD)
	dictionary(".C21.C9.C0.C1", "iPBinV6Address", valBCD)
	dictionary(".C21.C9.C0", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C21.C9.C0.C2", "iPTextV4Address", valSTRING)
	dictionary(".C21.C9.C0.C3", "iPTextV6Address", valSTRING)
	dictionary(".C21.C9.C1", "eTSIAddress", valBCD)
	dictionary(".C21.C12", "listOfTrafficVolumes", valSEQUENCE)
	dictionary(".C21.C12.U16.C2", "qosNegotiated", valBCD)
	dictionary(".C21.C12.U16.C3", "dataVolumeGPRSUplink", valINTEGER)
	dictionary(".C21.C12.U16.C4", "dataVolumeGPRSDownlink", valINTEGER)
	dictionary(".C21.C12.U16.C5", "changeCondition", valINTEGER)
	dictionary(".C21.C12.U16.C6", "changeTime", valBCD)
	dictionary(".C21.C12.U16.C8", "userLocationInformation", valBCD)
	dictionary(".C21.C19", "recordExtensions", valSEQUENCE)
	dictionary(".C21.C19.U16.U6", "identifier", valOBJID)
	dictionary(".C21.C19.U16.C1", "significance", valBOOLEAN)
	dictionary(".C21.C19.U16.C2", "information", valSET)
	dictionary(".C21.C19.U16.C2.C5", "userCategory", valINTEGER)
	dictionary(".C21.C19.U16.C2.C6", "ruleSpaceId", valSTRING)
	dictionary(".C21.C19.U16.C2.C2", "creditControlInfo", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C2.C6", "creditControlFailureReport", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C2.C6.C0", "requestType", valINTEGER)
	dictionary(".C21.C19.U16.C2.C2.C6.C1", "requestStatus", valINTEGER)
	dictionary(".C21.C19.U16.C2.C2.C6.C2", "resultCode", valINTEGER)
	dictionary(".C21.C19.U16.C2.C2.C6.C12", "ccRequestNumber", valINTEGER)
	dictionary(".C21.C19.U16.C2.C2.C7", "creditControlSessionId", valSTRING)
	dictionary(".C21.C19.U16.C2.C2.C8", "ccsRealm", valSTRING)
	dictionary(".C21.C19.U16.C2.C3", "policyControlInfo", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C3.C4", "policyControlFailureReport", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C3.C4.C0", "requestType", valINTEGER)
	dictionary(".C21.C19.U16.C2.C3.C4.C1", "requestStatus", valINTEGER)
	dictionary(".C21.C19.U16.C2.C3.C4.C2", "resultCode", valINTEGER)
	dictionary(".C21.C19.U16.C2.C3.C4.C5", "stopTime", valBCD)
	dictionary(".C21.C19.U16.C2.C3.C6", "pcsRealm", valSTRING)
	dictionary(".C21.C19.U16.C2.C3.C7", "policyControlSessionId", valSTRING)
	dictionary(".C21.C19.U16.C2.C7", "serviceContainers", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C7.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C2", "serviceIdentifier", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C3", "localSequenceNumber", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C4", "method", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C5", "inactivity", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C6", "resolution", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C7", "ccRequestNumber", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C8", "serviceSpecificUnits", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C9", "listOfURI", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C7.U16.C9.U16.C1", "count", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C9.U16.C2", "uri", valSTRING)
	dictionary(".C21.C19.U16.C2.C7.U16.C9.U16.C3", "uriIdentifier", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C9.U16.C4", "uriDataVolumeUplink", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C9.U16.C5", "uriDataVolumeDownlink", valINTEGER)
	dictionary(".C21.C19.U16.C2.C7.U16.C9.U16.C6", "listOfUriTimeStamps", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C8", "timeReports", valSEQUENCE)
	dictionary(".C21.C19.U16.C2.C8.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C21.C19.U16.C2.C8.U16.C2", "startTime", valBCD)
	dictionary(".C21.C19.U16.C2.C8.U16.C3", "endTime", valBCD)
	dictionary(".C21.C19.U16.C2.C8.U16.C4", "dataVolumeUplink", valINTEGER)
	dictionary(".C21.C19.U16.C2.C8.U16.C5", "dataVolumeDownlink", valINTEGER)
	dictionary(".C70", "egsnPDPRecord", valSET)
	dictionary(".C70.C0", "recordType", valINTEGER)
	dictionary(".C70.C3", "servedIMSI", valBCD)
	dictionary(".C70.C5", "chargingID", valINTEGER)
	dictionary(".C70.C7", "accessPointNameNI", valSTRING)
	dictionary(".C70.C8", "pdpType", valBCD)
	dictionary(".C70.C11", "dynamicAddressFlag", valBOOLEAN)
	dictionary(".C70.C13", "recordOpeningTime", valBCD)
	dictionary(".C70.C14", "duration", valINTEGER)
	dictionary(".C70.C15", "causeForRecClosing", valINTEGER)
	dictionary(".C70.C17", "recordSequenceNumber", valINTEGER)
	dictionary(".C70.C18", "nodeID", valSTRING)
	dictionary(".C70.C20", "localSequenceNumber", valINTEGER)
	dictionary(".C70.C21", "apnSelectionMode", valINTEGER)
	dictionary(".C70.C22", "servedMSISDN", valBCD)
	dictionary(".C70.C23", "chargingCharacteristics", valBCD)
	dictionary(".C70.C24", "chChSelectionMode", valINTEGER)
	dictionary(".C70.C25", "iMSsignalingContext", valNULL)
	dictionary(".C70.C27", "sgsnPLMNIdentifier", valBCD)
	dictionary(".C70.C29", "servedIMEISV", valBCD)
	dictionary(".C70.C30", "rATType", valINTEGER)
	dictionary(".C70.C31", "mSTimeZone", valBCD)
	dictionary(".C70.C32", "userLocationInformation", valBCD)
	dictionary(".C70.C4", "ggsnAddress", valSEQUENCE)
	dictionary(".C70.C4", "iPBinaryAddress", valCHOICE)
	dictionary(".C70.C4.C0", "iPBinV4Address", valBCD)
	dictionary(".C70.C4.C1", "iPBinV6Address", valBCD)
	dictionary(".C70.C4", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C70.C4.C2", "iPTextV4Address", valSTRING)
	dictionary(".C70.C4.C3", "iPTextV6Address", valSTRING)
	dictionary(".C70.C6", "sgsnAddress", valSEQUENCE)
	dictionary(".C70.C6", "iPBinaryAddress", valCHOICE)
	dictionary(".C70.C6.C0", "iPBinV4Address", valBCD)
	dictionary(".C70.C6.C1", "iPBinV6Address", valBCD)
	dictionary(".C70.C6", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C70.C6.C2", "iPTextV4Address", valSTRING)
	dictionary(".C70.C6.C3", "iPTextV6Address", valSTRING)
	dictionary(".C70.C9", "servedPDPAddress", valSEQUENCE)
	dictionary(".C70.C9.C0", "iPAddress", valSEQUENCE)
	dictionary(".C70.C9.C0", "iPBinaryAddress", valCHOICE)
	dictionary(".C70.C9.C0.C0", "iPBinV4Address", valBCD)
	dictionary(".C70.C9.C0.C1", "iPBinV6Address", valBCD)
	dictionary(".C70.C9.C0", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C70.C9.C0.C2", "iPTextV4Address", valSTRING)
	dictionary(".C70.C9.C0.C3", "iPTextV6Address", valSTRING)
	dictionary(".C70.C9.C1", "eTSIAddress", valBCD)
	dictionary(".C70.C12", "listOfTrafficVolumes", valSEQUENCE)
	dictionary(".C70.C12.U16.C2", "qosNegotiated", valBCD)
	dictionary(".C70.C12.U16.C3", "dataVolumeGPRSUplink", valINTEGER)
	dictionary(".C70.C12.U16.C4", "dataVolumeGPRSDownlink", valINTEGER)
	dictionary(".C70.C12.U16.C5", "changeCondition", valINTEGER)
	dictionary(".C70.C12.U16.C6", "changeTime", valBCD)
	dictionary(".C70.C12.U16.C8", "userLocationInformation", valBCD)
	dictionary(".C70.C19", "recordExtensions", valSEQUENCE)
	dictionary(".C70.C19.U16.U6", "identifier", valOBJID)
	dictionary(".C70.C19.U16.C1", "significance", valBOOLEAN)
	dictionary(".C70.C19.U16.C2", "information", valSET)
	dictionary(".C70.C19.U16.C2.C5", "userCategory", valINTEGER)
	dictionary(".C70.C19.U16.C2.C6", "ruleSpaceId", valSTRING)
	dictionary(".C70.C19.U16.C2.C2", "creditControlInfo", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C2.C6", "creditControlFailureReport", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C2.C6.C0", "requestType", valINTEGER)
	dictionary(".C70.C19.U16.C2.C2.C6.C1", "requestStatus", valINTEGER)
	dictionary(".C70.C19.U16.C2.C2.C6.C2", "resultCode", valINTEGER)
	dictionary(".C70.C19.U16.C2.C2.C6.C12", "ccRequestNumber", valINTEGER)
	dictionary(".C70.C19.U16.C2.C2.C7", "creditControlSessionId", valSTRING)
	dictionary(".C70.C19.U16.C2.C2.C8", "ccsRealm", valSTRING)
	dictionary(".C70.C19.U16.C2.C3", "policyControlInfo", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C3.C4", "policyControlFailureReport", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C3.C4.C0", "requestType", valINTEGER)
	dictionary(".C70.C19.U16.C2.C3.C4.C1", "requestStatus", valINTEGER)
	dictionary(".C70.C19.U16.C2.C3.C4.C2", "resultCode", valINTEGER)
	dictionary(".C70.C19.U16.C2.C3.C4.C5", "stopTime", valBCD)
	dictionary(".C70.C19.U16.C2.C3.C6", "pcsRealm", valSTRING)
	dictionary(".C70.C19.U16.C2.C3.C7", "policyControlSessionId", valSTRING)
	dictionary(".C70.C19.U16.C2.C7", "serviceContainers", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C7.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C2", "serviceIdentifier", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C3", "localSequenceNumber", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C4", "method", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C5", "inactivity", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C6", "resolution", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C7", "ccRequestNumber", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C8", "serviceSpecificUnits", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C9", "listOfURI", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C7.U16.C9.U16.C1", "count", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C9.U16.C2", "uri", valSTRING)
	dictionary(".C70.C19.U16.C2.C7.U16.C9.U16.C3", "uriIdentifier", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C9.U16.C4", "uriDataVolumeUplink", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C9.U16.C5", "uriDataVolumeDownlink", valINTEGER)
	dictionary(".C70.C19.U16.C2.C7.U16.C9.U16.C6", "listOfUriTimeStamps", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C8", "timeReports", valSEQUENCE)
	dictionary(".C70.C19.U16.C2.C8.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C70.C19.U16.C2.C8.U16.C2", "startTime", valBCD)
	dictionary(".C70.C19.U16.C2.C8.U16.C3", "endTime", valBCD)
	dictionary(".C70.C19.U16.C2.C8.U16.C4", "dataVolumeUplink", valINTEGER)
	dictionary(".C70.C19.U16.C2.C8.U16.C5", "dataVolumeDownlink", valINTEGER)
	dictionary(".C70.C28", "pSFurnishChargingInformation", valSEQUENCE)
	dictionary(".C70.C28.C1", "pSFreeFormatData", valSTRING)
	dictionary(".C70.C28.C2", "pSFFDAppendIndicator", valBOOLEAN)
	dictionary(".C70.C34", "listOfServiceData", valSEQUENCE)
	dictionary(".C70.C34.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C70.C34.U16.C3", "resultCode", valINTEGER)
	dictionary(".C70.C34.U16.C4", "localSequenceNumber", valINTEGER)
	dictionary(".C70.C34.U16.C5", "timeOfFirstUsage", valBCD)
	dictionary(".C70.C34.U16.C6", "timeOfLastUsage", valBCD)
	dictionary(".C70.C34.U16.C7", "timeUsage", valINTEGER)
	dictionary(".C70.C34.U16.C8", "serviceConditionChange", valBCD)
	dictionary(".C70.C34.U16.C9", "qoSInformationNeg", valBCD)
	dictionary(".C70.C34.U16.C10", "sgsn-Address", valSEQUENCE)
	dictionary(".C70.C34.U16.C10", "iPBinaryAddress", valCHOICE)
	dictionary(".C70.C34.U16.C10.C0", "iPBinV4Address", valBCD)
	dictionary(".C70.C34.U16.C10.C1", "iPBinV6Address", valBCD)
	dictionary(".C70.C34.U16.C10", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C70.C34.U16.C10.C2", "iPTextV4Address", valSTRING)
	dictionary(".C70.C34.U16.C10.C3", "iPTextV6Address", valSTRING)
	dictionary(".C70.C34.U16.C11", "sGSNPLMNIdentifier", valBCD)
	dictionary(".C70.C34.U16.C12", "datavolumeFBCUplink", valINTEGER)
	dictionary(".C70.C34.U16.C13", "datavolumeFBCDownlink", valINTEGER)
	dictionary(".C70.C34.U16.C14", "timeOfReport", valBCD)
	dictionary(".C70.C34.U16.C15", "rATType", valINTEGER)
	dictionary(".C70.C34.U16.C16", "failureHandlingContinue", valBOOLEAN)
	dictionary(".C70.C34.U16.C17", "serviceIdentifier", valINTEGER)
	dictionary(".C70.C34.U16.C18", "pSFurnishChargingInformation", valSEQUENCE)
	dictionary(".C70.C34.U16.C18.C1", "pSFreeFormatData", valSTRING)
	dictionary(".C70.C34.U16.C18.C2", "pSFFDAppendIndicator", valBOOLEAN)
	dictionary(".C70.C34.U16.C19", "aFRecordInformation", valSEQUENCE)
	dictionary(".C70.C34.U16.C20", "userLocationInformation", valBCD)
	dictionary(".C70.C34.U16.C21", "eventBasedChargingInformation", valSEQUENCE)
	dictionary(".C70.C34.U16.C21.C1", "numberOfEvents", valINTEGER)
	dictionary(".C70.C34.U16.C21.C2", "eventTimeStamps", valSEQUENCE)
	dictionary(".C78", "sGWRecord", valSET)
	dictionary(".C78.C0", "recordType", valINTEGER)
	dictionary(".C78.C3", "servedIMSI", valBCD)
	dictionary(".C78.C5", "chargingID", valINTEGER)
	dictionary(".C78.C7", "accessPointNameNI", valSTRING)
	dictionary(".C78.C8", "pdpPDNType", valBCD)
	dictionary(".C78.C13", "recordOpeningTime", valBCD)
	dictionary(".C78.C14", "duration", valINTEGER)
	dictionary(".C78.C15", "causeForRecClosing", valINTEGER)
	dictionary(".C78.C17", "recordSequenceNumber", valINTEGER)
	dictionary(".C78.C18", "nodeID", valSTRING)
	dictionary(".C78.C20", "localSequenceNumber", valINTEGER)
	dictionary(".C78.C22", "servedMSISDN", valBCD)
	dictionary(".C78.C23", "chargingCharacteristics", valBCD)
	dictionary(".C78.C27", "servingNodePLMNIdentifier", valBCD)
	dictionary(".C78.C29", "servedIMEISV", valBCD)
	dictionary(".C78.C30", "rATType", valINTEGER)
	dictionary(".C78.C31", "mSTimeZone", valBCD)
	dictionary(".C78.C34", "sGWChange", valBOOLEAN)
	dictionary(".C78.C37", "p_GWPLMNIdentifier", valBCD)
	dictionary(".C78.C40", "pDNConnectionID", valINTEGER)
	dictionary(".C78.C4", "s_GWAddress", valSEQUENCE)
	dictionary(".C78.C4", "iPBinaryAddress", valCHOICE)
	dictionary(".C78.C4.C0", "iPBinV4Address", valBCD)
	dictionary(".C78.C4.C1", "iPBinV6Address", valBCD)
	dictionary(".C78.C4", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C78.C4.C2", "iPTextV4Address", valSTRING)
	dictionary(".C78.C4.C3", "iPTextV6Address", valSTRING)
	dictionary(".C78.C6", "servingNodeAddress", valSEQUENCE)
	dictionary(".C78.C6", "iPBinaryAddress", valCHOICE)
	dictionary(".C78.C6.C0", "iPBinV4Address", valBCD)
	dictionary(".C78.C6.C1", "iPBinV6Address", valBCD)
	dictionary(".C78.C6", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C78.C6.C2", "iPTextV4Address", valSTRING)
	dictionary(".C78.C6.C3", "iPTextV6Address", valSTRING)
	dictionary(".C78.C9", "servedPDPPDNAddress", valSEQUENCE)
	dictionary(".C78.C9.C0", "iPAddress", valSEQUENCE)
	dictionary(".C78.C9.C0", "iPBinaryAddress", valCHOICE)
	dictionary(".C78.C9.C0.C0", "iPBinV4Address", valBCD)
	dictionary(".C78.C9.C0.C1", "iPBinV6Address", valBCD)
	dictionary(".C78.C9.C0", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C78.C9.C0.C2", "iPTextV4Address", valSTRING)
	dictionary(".C78.C9.C0.C3", "iPTextV6Address", valSTRING)
	dictionary(".C78.C9.C1", "eTSIAddress", valBCD)
	dictionary(".C78.C12", "listOfTrafficVolumes", valSEQUENCE)
	dictionary(".C78.C12.U16.C3", "dataVolumeGPRSUplink", valINTEGER)
	dictionary(".C78.C12.U16.C4", "dataVolumeGPRSDownlink", valINTEGER)
	dictionary(".C78.C12.U16.C5", "changeCondition", valINTEGER)
	dictionary(".C78.C12.U16.C6", "changeTime", valBCD)
	dictionary(".C78.C12.U16.C8", "userLocationInformation", valBCD)
	dictionary(".C78.C12.U16.C9", "ePCQoSInformation", valSEQUENCE)
	dictionary(".C78.C12.U16.C9.C1", "qCI", valINTEGER)
	dictionary(".C78.C12.U16.C9.C2", "maxRequestedBandwithUL", valINTEGER)
	dictionary(".C78.C12.U16.C9.C3", "maxRequestedBandwithDL", valINTEGER)
	dictionary(".C78.C12.U16.C9.C4", "guaranteedBitrateUL", valINTEGER)
	dictionary(".C78.C12.U16.C9.C5", "guaranteedBitrateDL", valINTEGER)
	dictionary(".C78.C12.U16.C9.C6", "aRP", valINTEGER)
	dictionary(".C78.C35", "servingNodeType", valSEQUENCE)
	dictionary(".C78.C36", "p_GWAddressUsed", valSEQUENCE)
	dictionary(".C78.C36", "iPBinaryAddress", valCHOICE)
	dictionary(".C78.C36.C0", "iPBinV4Address", valBCD)
	dictionary(".C78.C36.C1", "iPBinV6Address", valBCD)
	dictionary(".C78.C36", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C78.C36.C2", "iPTextV4Address", valSTRING)
	dictionary(".C78.C36.C3", "iPTextV6Address", valSTRING)
	dictionary(".C78.C43", "servedPDPPDNAddressExt", valSEQUENCE)
	dictionary(".C78.C43.C0", "iPAddress", valSEQUENCE)
	dictionary(".C78.C43.C0", "iPBinaryAddress", valCHOICE)
	dictionary(".C78.C43.C0.C0", "iPBinV4Address", valBCD)
	dictionary(".C78.C43.C0.C1", "iPBinV6Address", valBCD)
	dictionary(".C78.C43.C0", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C78.C43.C0.C2", "iPTextV4Address", valSTRING)
	dictionary(".C78.C43.C0.C3", "iPTextV6Address", valSTRING)
	dictionary(".C78.C43.C1", "eTSIAddress", valBCD)
	dictionary(".C79", "pgwRecord", valSET)
	dictionary(".C79.C0", "recordType", valINTEGER)
	dictionary(".C79.C3", "servedIMSI", valBCD)
	dictionary(".C79.C5", "chargingID", valINTEGER)
	dictionary(".C79.C7", "accessPointNameNI", valSTRING)
	dictionary(".C79.C8", "pdpPDNType", valBCD)
	dictionary(".C79.C11", "dynamicAddressFlag", valBOOLEAN)
	dictionary(".C79.C13", "recordOpeningTime", valBCD)
	dictionary(".C79.C14", "duration", valINTEGER)
	dictionary(".C79.C15", "causeForRecClosing", valINTEGER)
	dictionary(".C79.C17", "recordSequenceNumber", valINTEGER)
	dictionary(".C79.C18", "nodeID", valSTRING)
	dictionary(".C79.C20", "localSequenceNumber", valINTEGER)
	dictionary(".C79.C21", "apnSelectionMode", valINTEGER)
	dictionary(".C79.C22", "servedMSISDN", valBCD)
	dictionary(".C79.C23", "chargingCharacteristics", valBCD)
	dictionary(".C79.C24", "chChSelectionMode", valINTEGER)
	dictionary(".C79.C25", "iMSsignalingContext", valNULL)
	dictionary(".C79.C27", "servingNodePLMNIdentifier", valBCD)
	dictionary(".C79.C29", "servedIMEISV", valBCD)
	dictionary(".C79.C30", "rATType", valINTEGER)
	dictionary(".C79.C31", "mSTimeZone", valBCD)
	dictionary(".C79.C32", "userLocationInformation", valBCD)
	dictionary(".C79.C37", "p_GWPLMNIdentifier", valBCD)
	dictionary(".C79.C38", "startTime", valBCD)
	dictionary(".C79.C39", "stopTime", valBCD)
	dictionary(".C79.C41", "pDNConnectionID", valINTEGER)
	dictionary(".C79.C4", "p_GWAddress", valSEQUENCE)
	dictionary(".C79.C4", "iPBinaryAddress", valCHOICE)
	dictionary(".C79.C4.C0", "iPBinV4Address", valBCD)
	dictionary(".C79.C4.C1", "iPBinV6Address", valBCD)
	dictionary(".C79.C4", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C79.C4.C2", "iPTextV4Address", valSTRING)
	dictionary(".C79.C4.C3", "iPTextV6Address", valSTRING)
	dictionary(".C79.C6", "servingNodeAddress", valSEQUENCE)
	dictionary(".C79.C6", "iPBinaryAddress", valCHOICE)
	dictionary(".C79.C6.C0", "iPBinV4Address", valBCD)
	dictionary(".C79.C6.C1", "iPBinV6Address", valBCD)
	dictionary(".C79.C6", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C79.C6.C2", "iPTextV4Address", valSTRING)
	dictionary(".C79.C6.C3", "iPTextV6Address", valSTRING)
	dictionary(".C79.C9", "servedPDPPDNAddress", valSEQUENCE)
	dictionary(".C79.C9.C0", "iPAddress", valSEQUENCE)
	dictionary(".C79.C9.C0", "iPBinaryAddress", valCHOICE)
	dictionary(".C79.C9.C0.C0", "iPBinV4Address", valBCD)
	dictionary(".C79.C9.C0.C1", "iPBinV6Address", valBCD)
	dictionary(".C79.C9.C0", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C79.C9.C0.C2", "iPTextV4Address", valSTRING)
	dictionary(".C79.C9.C0.C3", "iPTextV6Address", valSTRING)
	dictionary(".C79.C9.C1", "eTSIAddress", valBCD)
	dictionary(".C79.C19", "recordExtensions", valSEQUENCE)
	dictionary(".C79.C19.U16.U6", "identifier", valOBJID)
	dictionary(".C79.C19.U16.C1", "significance", valBOOLEAN)
	dictionary(".C79.C19.U16.C2", "information", valSET)
	dictionary(".C79.C19.U16.C2.C5", "userCategory", valINTEGER)
	dictionary(".C79.C19.U16.C2.C6", "ruleSpaceId", valSTRING)
	dictionary(".C79.C19.U16.C2.C2", "creditControlInfo", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C2.C6", "creditControlFailureReport", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C2.C6.C0", "requestType", valINTEGER)
	dictionary(".C79.C19.U16.C2.C2.C6.C1", "requestStatus", valINTEGER)
	dictionary(".C79.C19.U16.C2.C2.C6.C2", "resultCode", valINTEGER)
	dictionary(".C79.C19.U16.C2.C2.C6.C12", "ccRequestNumber", valINTEGER)
	dictionary(".C79.C19.U16.C2.C2.C7", "creditControlSessionId", valSTRING)
	dictionary(".C79.C19.U16.C2.C2.C8", "ccsRealm", valSTRING)
	dictionary(".C79.C19.U16.C2.C3", "policyControlInfo", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C3.C4", "policyControlFailureReport", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C3.C4.C0", "requestType", valINTEGER)
	dictionary(".C79.C19.U16.C2.C3.C4.C1", "requestStatus", valINTEGER)
	dictionary(".C79.C19.U16.C2.C3.C4.C2", "resultCode", valINTEGER)
	dictionary(".C79.C19.U16.C2.C3.C4.C5", "stopTime", valBCD)
	dictionary(".C79.C19.U16.C2.C3.C6", "pcsRealm", valSTRING)
	dictionary(".C79.C19.U16.C2.C3.C7", "policyControlSessionId", valSTRING)
	dictionary(".C79.C19.U16.C2.C7", "serviceContainers", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C7.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C2", "serviceIdentifier", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C3", "localSequenceNumber", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C4", "method", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C5", "inactivity", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C6", "resolution", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C7", "ccRequestNumber", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C8", "serviceSpecificUnits", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C9", "listOfURI", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C7.U16.C9.U16.C1", "count", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C9.U16.C2", "uri", valSTRING)
	dictionary(".C79.C19.U16.C2.C7.U16.C9.U16.C3", "uriIdentifier", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C9.U16.C4", "uriDataVolumeUplink", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C9.U16.C5", "uriDataVolumeDownlink", valINTEGER)
	dictionary(".C79.C19.U16.C2.C7.U16.C9.U16.C6", "listOfUriTimeStamps", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C8", "timeReports", valSEQUENCE)
	dictionary(".C79.C19.U16.C2.C8.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C79.C19.U16.C2.C8.U16.C2", "startTime", valBCD)
	dictionary(".C79.C19.U16.C2.C8.U16.C3", "endTime", valBCD)
	dictionary(".C79.C19.U16.C2.C8.U16.C4", "dataVolumeUplink", valINTEGER)
	dictionary(".C79.C19.U16.C2.C8.U16.C5", "dataVolumeDownlink", valINTEGER)
	dictionary(".C79.C28", "pSFurnishChargingInformation", valSEQUENCE)
	dictionary(".C79.C28.C1", "pSFreeFormatData", valSTRING)
	dictionary(".C79.C28.C2", "pSFFDAppendIndicator", valBOOLEAN)
	dictionary(".C79.C34", "listOfServiceData", valSEQUENCE)
	dictionary(".C79.C34.U16.C1", "ratingGroup", valINTEGER)
	dictionary(".C79.C34.U16.C3", "resultCode", valINTEGER)
	dictionary(".C79.C34.U16.C4", "localSequenceNumber", valINTEGER)
	dictionary(".C79.C34.U16.C5", "timeOfFirstUsage", valBCD)
	dictionary(".C79.C34.U16.C6", "timeOfLastUsage", valBCD)
	dictionary(".C79.C34.U16.C7", "timeUsage", valINTEGER)
	dictionary(".C79.C34.U16.C8", "serviceConditionChange", valBCD)
	dictionary(".C79.C34.U16.C9", "qoSInformationNeg", valSEQUENCE)
	dictionary(".C79.C34.U16.C9.C1", "qCI", valINTEGER)
	dictionary(".C79.C34.U16.C9.C2", "maxRequestedBandwithUL", valINTEGER)
	dictionary(".C79.C34.U16.C9.C3", "maxRequestedBandwithDL", valINTEGER)
	dictionary(".C79.C34.U16.C9.C4", "guaranteedBitrateUL", valINTEGER)
	dictionary(".C79.C34.U16.C9.C5", "guaranteedBitrateDL", valINTEGER)
	dictionary(".C79.C34.U16.C9.C6", "aRP", valINTEGER)
	dictionary(".C79.C34.U16.C10", "sgsn-Address", valSEQUENCE)
	dictionary(".C79.C34.U16.C10", "iPBinaryAddress", valCHOICE)
	dictionary(".C79.C34.U16.C10.C0", "iPBinV4Address", valBCD)
	dictionary(".C79.C34.U16.C10.C1", "iPBinV6Address", valBCD)
	dictionary(".C79.C34.U16.C10", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C79.C34.U16.C10.C2", "iPTextV4Address", valSTRING)
	dictionary(".C79.C34.U16.C10.C3", "iPTextV6Address", valSTRING)
	dictionary(".C79.C34.U16.C11", "sGSNPLMNIdentifier", valBCD)
	dictionary(".C79.C34.U16.C12", "datavolumeFBCUplink", valINTEGER)
	dictionary(".C79.C34.U16.C13", "datavolumeFBCDownlink", valINTEGER)
	dictionary(".C79.C34.U16.C14", "timeOfReport", valBCD)
	dictionary(".C79.C34.U16.C15", "rATType", valINTEGER)
	dictionary(".C79.C34.U16.C16", "failureHandlingContinue", valBOOLEAN)
	dictionary(".C79.C34.U16.C17", "serviceIdentifier", valINTEGER)
	dictionary(".C79.C34.U16.C18", "pSFurnishChargingInformation", valSEQUENCE)
	dictionary(".C79.C34.U16.C18.C1", "pSFreeFormatData", valSTRING)
	dictionary(".C79.C34.U16.C18.C2", "pSFFDAppendIndicator", valBOOLEAN)
	dictionary(".C79.C34.U16.C19", "aFRecordInformation", valSEQUENCE)
	dictionary(".C79.C34.U16.C19.U16.C1", "aFChargingIdentifier", valSTRING)
	dictionary(".C79.C34.U16.C20", "userLocationInformation", valBCD)
	dictionary(".C79.C34.U16.C21", "eventBasedChargingInformation", valSEQUENCE)
	dictionary(".C79.C34.U16.C21.C1", "numberOfEvents", valINTEGER)
	dictionary(".C79.C34.U16.C21.C2", "eventTimeStamps", valSEQUENCE)
	dictionary(".C79.C35", "servingNodeType", valSEQUENCE)
	dictionary(".C79.C45", "servedPDPPDNAddressExt", valSEQUENCE)
	dictionary(".C79.C45.C0", "iPAddress", valSEQUENCE)
	dictionary(".C79.C45.C0", "iPBinaryAddress", valCHOICE)
	dictionary(".C79.C45.C0.C0", "iPBinV4Address", valBCD)
	dictionary(".C79.C45.C0.C1", "iPBinV6Address", valBCD)
	dictionary(".C79.C45.C0", "iPTextRepresentedAddress", valCHOICE)
	dictionary(".C79.C45.C0.C2", "iPTextV4Address", valSTRING)
	dictionary(".C79.C45.C0.C3", "iPTextV6Address", valSTRING)
	dictionary(".C79.C45.C1", "eTSIAddress", valBCD)

	return true
}
