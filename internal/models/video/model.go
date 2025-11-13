package video

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID                 string         `json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at"`
	MuxUploadID        *string        `json:"mux_upload_id"`
	MuxAssetID         *string        `json:"mux_asset_id"`
	MuxPlaybackID      *string        `json:"mux_playback_id"`
	State              string         `json:"state"`
	Status             *string        `json:"status"`
	Duration           *float32       `json:"duration"`
	AspectRatio        *string        `json:"aspect_ratio"`
	Height             *int           `json:"height"`
	Width              *int           `json:"width"`
	AssetCreatedAt     *time.Time     `json:"asset_created_at"`
	ResolutionTier     *string        `json:"resolution_tier"`
	MaxStoredFrameRate *string        `json:"max_stored_frame_rate"`
	IngestType         *string        `json:"ingest_type"`
	Passthrough        *string        `json:"passthrough"`
	OwnerID            *string        `json:"owner_id"`
	OwnerType          *string        `json:"owner_type"`
}
