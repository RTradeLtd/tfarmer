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
func NewFarmer(db *gorm.DB) *Farmer {
	return &Farmer{
		UM: models.NewUploadManager(db),
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
func (f *Farmer) AverageUploadSize() (float64, error) {
	uploads, err := f.UM.GetUploads()
	if err != nil {
		return 0, err
	}
	var totalSizeInBytes int
	for _, v := range uploads {
		stats, err := f.ipfs.Stat(v.Hash)
		if err != nil {
			return 0, err
		}
		totalSizeInBytes = totalSizeInBytes + stats.CumulativeSize
	}
	totalSizeInGigaBytes := float64(totalSizeInBytes) / float64(datasize.GB.Bytes())
	averageSizeInGigaBytes := totalSizeInGigaBytes / float64(len(uploads))
	return averageSizeInGigaBytes, nil
}
