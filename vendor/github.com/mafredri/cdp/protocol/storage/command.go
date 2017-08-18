// Code generated by cdpgen. DO NOT EDIT.

package storage

// ClearDataForOriginArgs represents the arguments for ClearDataForOrigin in the Storage domain.
type ClearDataForOriginArgs struct {
	Origin       string `json:"origin"`       // Security origin.
	StorageTypes string `json:"storageTypes"` // Comma separated origin names.
}

// NewClearDataForOriginArgs initializes ClearDataForOriginArgs with the required arguments.
func NewClearDataForOriginArgs(origin string, storageTypes string) *ClearDataForOriginArgs {
	args := new(ClearDataForOriginArgs)
	args.Origin = origin
	args.StorageTypes = storageTypes
	return args
}

// GetUsageAndQuotaArgs represents the arguments for GetUsageAndQuota in the Storage domain.
type GetUsageAndQuotaArgs struct {
	Origin string `json:"origin"` // Security origin.
}

// NewGetUsageAndQuotaArgs initializes GetUsageAndQuotaArgs with the required arguments.
func NewGetUsageAndQuotaArgs(origin string) *GetUsageAndQuotaArgs {
	args := new(GetUsageAndQuotaArgs)
	args.Origin = origin
	return args
}

// GetUsageAndQuotaReply represents the return values for GetUsageAndQuota in the Storage domain.
type GetUsageAndQuotaReply struct {
	Usage          float64        `json:"usage"`          // Storage usage (bytes).
	Quota          float64        `json:"quota"`          // Storage quota (bytes).
	UsageBreakdown []UsageForType `json:"usageBreakdown"` // Storage usage per type (bytes).
}

// TrackCacheStorageForOriginArgs represents the arguments for TrackCacheStorageForOrigin in the Storage domain.
type TrackCacheStorageForOriginArgs struct {
	Origin string `json:"origin"` // Security origin.
}

// NewTrackCacheStorageForOriginArgs initializes TrackCacheStorageForOriginArgs with the required arguments.
func NewTrackCacheStorageForOriginArgs(origin string) *TrackCacheStorageForOriginArgs {
	args := new(TrackCacheStorageForOriginArgs)
	args.Origin = origin
	return args
}

// UntrackCacheStorageForOriginArgs represents the arguments for UntrackCacheStorageForOrigin in the Storage domain.
type UntrackCacheStorageForOriginArgs struct {
	Origin string `json:"origin"` // Security origin.
}

// NewUntrackCacheStorageForOriginArgs initializes UntrackCacheStorageForOriginArgs with the required arguments.
func NewUntrackCacheStorageForOriginArgs(origin string) *UntrackCacheStorageForOriginArgs {
	args := new(UntrackCacheStorageForOriginArgs)
	args.Origin = origin
	return args
}
