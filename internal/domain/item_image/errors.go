package itemimage

import "errors"

var ErrPhotoNotOwned = errors.New("cannot add photo to item: item does not belong to the user")