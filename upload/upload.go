package upload

import (
	"github.com/RTradeLtd/database/models"
	"github.com/RTradeLtd/gorm"
	"github.com/RTradeLtd/rtfs"
	"github.com/c2h5oh/datasize"
)

// Farmer is used to gather upload information
type Farmer struct {
	UM   *models.UploadManager
	ipfs rtfs.Manager
}

// NewFarmer is used to instantiate our upload farmer
func NewFarmer(db *gorm.DB, ipfs rtfs.Manager) *Farmer {
	return &Farmer{
		UM:   models.NewUploadManager(db),
		ipfs: ipfs,
	}
}

// TotalUploadsCount is used to get the total number of uploads
// allows control of whether or not we count unique uploads.
//
// when counting by unique uploads, if the same content hash is
// uploaded by multiple users, then we only count it once.
func (f *Farmer) TotalUploadsCount(unique bool) (int, error) {
	uploads, err := f.UM.GetUploads()
	if err != nil {
		return 0, err
	}
	var numberOfUploads int
	if unique {
		found := make(map[string]bool)
		for _, v := range uploads {
			if !found[v.Hash] {
				found[v.Hash] = true
			}
		}
		// len() on a map type gives the number of keys
		numberOfUploads = len(found)
	} else {
		numberOfUploads = len(uploads)
	}
	return numberOfUploads, nil
}

// AverageUploadSize is used to get the average size of uploads
func (f *Farmer) AverageUploadSize(unique bool) (float64, error) {
	uploads, err := f.UM.GetUploads()
	if err != nil {
		return 0, err
	}
	var (
		totalSizeInBytes int
		numUploads       int
	)
	if unique {
		found := make(map[string]bool)
		for _, v := range uploads {
			if !found[v.Hash] {
				found[v.Hash] = true
				stats, err := f.ipfs.Stat(v.Hash)
				if err != nil {
					return 0, err
				}
				totalSizeInBytes = totalSizeInBytes + stats.CumulativeSize
			}
		}
		numUploads = len(found)
	} else {
		for _, v := range uploads {
			stats, err := f.ipfs.Stat(v.Hash)
			if err != nil {
				return 0, err
			}
			totalSizeInBytes = totalSizeInBytes + stats.CumulativeSize
		}
		numUploads = len(uploads)
	}
	totalSizeInGigaBytes := float64(totalSizeInBytes) / float64(datasize.GB.Bytes())
	averageSizeInGigaBytes := totalSizeInGigaBytes / float64(numUploads)
	return averageSizeInGigaBytes, nil
}

// NumberOfUploads is used to retrieve number of uploads by type
func (f *Farmer) NumberOfUploads() (int, error) {
	uploads := []models.Upload{}
	if err := f.UM.DB.Find(&uploads).Error; err != nil {
		return 0, err
	}
	return len(uploads), nil
}
