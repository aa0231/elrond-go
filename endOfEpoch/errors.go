package endOfEpoch

import "errors"

// ErrNilArgsNewMetaEndOfEpochTrigger signals that nil arguments were provided
var ErrNilArgsNewMetaEndOfEpochTrigger = errors.New("nil arguments for meta end of epoch trigger")

// ErrNilRounder signals that nil round was provided
var ErrNilRounder = errors.New("nil rounder")

// ErrNilEndOfEpochSettings signals that nil end of epoch settings has been provided
var ErrNilEndOfEpochSettings = errors.New("nil end of epoch settings")

// ErrInvalidSettingsForEndOfEpochTrigger signals that settings for end of epoch trigger are invalid
var ErrInvalidSettingsForEndOfEpochTrigger = errors.New("invalid end of epoch trigger settings")

// ErrNilSyncTimer signals that sync timer is nil
var ErrNilSyncTimer = errors.New("nil sync timer")

// ErrNilArgsNewShardEndOfEpochTrigger signals that nil arguments for shard epoch trigger has been provided
var ErrNilArgsNewShardEndOfEpochTrigger = errors.New("nil arguments for shard end of epoch trigger")

// ErrNotEnoughRoundsBetweenEpochs signals that not enough rounds has passed since last epoch start
var ErrNotEnoughRoundsBetweenEpochs = errors.New("tried to force end of epoch before passing of enough rounds")

// ErrForceEndOfEpochCanBeCalledOnNewRound signals that force end of epoch was called on wrong round
var ErrForceEndOfEpochCanBeCalledOnNewRound = errors.New("invalid time to call force end of epoch, possible only on new round")

// ErrSavedRoundIsHigherThanSaved signals that input round was wrong
var ErrSavedRoundIsHigherThanSaved = errors.New("saved round is higher than input round")

// ErrWrongTypeAssertion signals wrong type assertion
var ErrWrongTypeAssertion = errors.New("wrong type assertion")

// ErrNilMarshalizer signals that nil marshalizer has been provided
var ErrNilMarshalizer = errors.New("nil marshalizer")

// ErrNilStorage signals that nil storage has been provided
var ErrNilStorage = errors.New("nil storage")

// ErrNilHeaderHandler signal that a nil header handler has been provided
var ErrNilHeaderHandler = errors.New("nil header handler")

// ErrNilArgsPendingMiniblocks signals that nil argument was passed
var ErrNilArgsPendingMiniblocks = errors.New("nil arguments for pending miniblock object")

// ErrMetaHdrNotFound signals that metaheader was not found
var ErrMetaHdrNotFound = errors.New("meta header not found")

// ErrNilHasher signals that nil hasher has been provided
var ErrNilHasher = errors.New("nil hasher")

// ErrNilHeaderValidator signals that nil header validator has been provided
var ErrNilHeaderValidator = errors.New("nil header validator")

// ErrNilDataPoolsHolder signals that nil data pools holder has been provided
var ErrNilDataPoolsHolder = errors.New("nil data pools holder")

// ErrNilStorageService signals that nil storage service has been provided
var ErrNilStorageService = errors.New("nil storage service")

// ErrNilRequestHandler signals that nil request handler has been provided
var ErrNilRequestHandler = errors.New("nil request handler")

// ErrNilMetaBlocksPool signals that nil metablock pools holder has been provided
var ErrNilMetaBlocksPool = errors.New("nil metablocks pool")

// ErrNilHeaderNoncesPool signals that nil header nonces pool has been provided
var ErrNilHeaderNoncesPool = errors.New("nil header nonces pool")

// ErrNilUint64Converter signals that nil uint64 converter has been provided
var ErrNilUint64Converter = errors.New("nil uint64 converter")

// ErrNilMetaHdrStorage signals that nil meta header storage has been provided
var ErrNilMetaHdrStorage = errors.New("nil meta header storage")

// ErrNilMetaNonceHashStorage signals that nil meta header nonce hash storage has been provided
var ErrNilMetaNonceHashStorage = errors.New("nil meta nonce hash storage")
