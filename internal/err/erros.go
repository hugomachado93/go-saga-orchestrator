package err

import "errors"

var ErroStateOrEventNotFound = errors.New("State or event not found")

var ErroFailedFindApiKey = errors.New("failed to find x-api-key")
