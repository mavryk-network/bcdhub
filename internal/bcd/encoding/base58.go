package encoding

import (
	"encoding/hex"
	"errors"

	"github.com/ebellocchia/go-base58"
)

var base58Enc = base58.New(base58.AlphabetBitcoin)

type base58Encoding struct {
	EncodedPrefix []byte
	EncodedLength int
	DecodedPrefix []byte
	DecodedLength int
	DataType      string
}

// Bytes prefixes
const (
	PrefixBlockHash                   = "B"
	PrefixOperationHash               = "o"
	PrefixOperationListHash           = "Lo"
	PrefixOperationListListHash       = "LLo"
	PrefixProtocolHash                = "P"
	PrefixContextHash                 = "Co"
	PrefixPublicKeyTZ1                = "mv1"
	PrefixPublicKeyTZ2                = "mv2"
	PrefixPublicKeyTZ3                = "mv3"
	PrefixPublicKeyTZ4                = "mv4"
	PrefixPublicKeyKT1                = "KT1"
	PrefixPublicKeyTxr1               = "txr1"
	PrefixOriginatedSmartRollup       = "sr1"
	PrefixRollupInboxHash             = "txi"
	PrefixRollupMessageHash           = "txm"
	PrefixRollupCommitmentHash        = "txc"
	PrefixRollupMessageResultHash     = "txmr"
	PrefixRollupMessageResultListHash = "txM"
	PrefixRollupWithdrawListHash      = "txw"
	PrefixScriptExpr                  = "expr"
	PrefixED25519Seed                 = "edsk"
	PrefixED25519PublicKey            = "edpk"
	PrefixSecp256k1SecretKey          = "spsk"
	PrefixP256SecretKey               = "p2sk"
	PrefixSecp256k1PublicKey          = "sppk"
	PrefixP256PublicKey               = "p2pk"
	PrefixSecp256k1Scalar             = "SSp"
	PrefixSecp256k1Element            = "GSp"
	PrefixED25519SecretKey            = "edsk"
	PrefixED25519Signature            = "edsig"
	PrefixSecp256k1Signature          = "spsig"
	PrefixP256Signature               = "p2sig"
	PrefixGenericSignature            = "sig"
	PrefixChainID                     = "Net"
	PrefixCryptoBoxPublicKeyHash      = "id"
	PrefixED25519EncryptedSeed        = "edesk"
	PrefixSecp256k1EncryptedSecretKey = "spesk"
	PrefixP256EncryptedSecretKey      = "p2esk"
	PrefixBakerHash                   = "SG1"
	PrefixSmartRollupCommitment       = "src1"
	PrefixSmartRollupState            = "srs1"
	PrefixSmartRollupInbox            = "srib1"
	PrefixSmartRollupMerkelizedTree   = "srib2"
)

var base58Encodings = []base58Encoding{
	{[]byte(PrefixBlockHash), 51, []byte{1, 52}, 32, "block hash"},
	{[]byte(PrefixOperationHash), 51, []byte{5, 116}, 32, "operation hash"},
	{[]byte(PrefixOperationListHash), 52, []byte{133, 233}, 32, "operation list hash"},
	{[]byte(PrefixOperationListListHash), 53, []byte{29, 159, 109}, 32, "operation list list hash"},
	{[]byte(PrefixProtocolHash), 51, []byte{2, 170}, 32, "protocol hash"},
	{[]byte(PrefixContextHash), 52, []byte{79, 199}, 32, "context hash"},

	{[]byte(PrefixPublicKeyTZ1), 36, []byte{5, 186, 196}, 20, "ed25519 public key hash"},
	{[]byte(PrefixPublicKeyTZ2), 36, []byte{5, 186, 199}, 20, "secp256k1 public key hash"},
	{[]byte(PrefixPublicKeyTZ3), 36, []byte{5, 186, 201}, 20, "p256 public key hash"},
	{[]byte(PrefixPublicKeyTZ4), 36, []byte{5, 186, 204}, 20, "BLS-MinPk"},
	{[]byte(PrefixPublicKeyKT1), 36, []byte{2, 90, 121}, 20, "Originated address"},

	{[]byte(PrefixPublicKeyTxr1), 37, []byte{1, 128, 120, 31}, 20, "Rollup address"},
	{[]byte(PrefixOriginatedSmartRollup), 37, []byte{6, 124, 117}, 20, "Smart rollup"},
	{[]byte(PrefixRollupInboxHash), 53, []byte{79, 148, 196}, 32, "Inbox hash"},
	{[]byte(PrefixRollupMessageHash), 53, []byte{79, 149, 030}, 32, "Message hash"},
	{[]byte(PrefixRollupCommitmentHash), 53, []byte{79, 148, 017}, 32, "Commitment hash"},
	{[]byte(PrefixRollupMessageResultHash), 54, []byte{18, 7, 206, 87}, 32, "Message result hash"},
	{[]byte(PrefixRollupMessageResultListHash), 53, []byte{79, 146, 82}, 32, "Message result list hash"},
	{[]byte(PrefixRollupWithdrawListHash), 53, []byte{79, 150, 72}, 32, "Withdraw list hash"},

	{[]byte(PrefixScriptExpr), 54, []byte{13, 44, 64, 27}, 32, "script expression"},
	{[]byte(PrefixED25519Seed), 54, []byte{13, 15, 58, 7}, 32, "ed25519 seed"},
	{[]byte(PrefixED25519PublicKey), 54, []byte{13, 15, 37, 217}, 32, "ed25519 public key"},
	{[]byte(PrefixSecp256k1SecretKey), 54, []byte{17, 162, 224, 201}, 32, "secp256k1 secret key"},
	{[]byte(PrefixP256SecretKey), 54, []byte{16, 81, 238, 189}, 32, "p256 secret key"},

	{[]byte(PrefixSecp256k1PublicKey), 55, []byte{3, 254, 226, 86}, 33, "secp256k1 public key"},
	{[]byte(PrefixP256PublicKey), 55, []byte{3, 178, 139, 127}, 33, "p256 public key"},
	{[]byte(PrefixSecp256k1Scalar), 53, []byte{38, 248, 136}, 33, "secp256k1 scalar"},
	{[]byte(PrefixSecp256k1Element), 53, []byte{5, 92, 0}, 33, "secp256k1 element"},

	{[]byte(PrefixED25519SecretKey), 98, []byte{43, 246, 78, 7}, 64, "ed25519 secret key"},
	{[]byte(PrefixED25519Signature), 99, []byte{9, 245, 205, 134, 18}, 64, "ed25519 signature"},
	{[]byte(PrefixSecp256k1Signature), 99, []byte{13, 115, 101, 19, 63}, 64, "secp256k1 signature"},
	{[]byte(PrefixP256Signature), 98, []byte{54, 240, 44, 52}, 64, "p256 signature"},
	{[]byte(PrefixGenericSignature), 96, []byte{4, 130, 43}, 64, "generic signature"},

	{[]byte(PrefixChainID), 15, []byte{87, 82, 0}, 4, "chain id"},

	{[]byte(PrefixCryptoBoxPublicKeyHash), 30, []byte{153, 103}, 16, "cryptobox public key hash"},

	{[]byte(PrefixED25519EncryptedSeed), 88, []byte{7, 90, 60, 179, 41}, 56, "ed25519 encrypted seed"},
	{[]byte(PrefixSecp256k1EncryptedSecretKey), 88, []byte{9, 237, 241, 174, 150}, 56, "secp256k1 encrypted secret key"},
	{[]byte(PrefixP256EncryptedSecretKey), 88, []byte{9, 48, 57, 115, 171}, 56, "p256_encrypted_secret_key"},

	{[]byte(PrefixBakerHash), 36, []byte{3, 56, 226}, 20, "baker hash"},
	{[]byte(PrefixSmartRollupCommitment), 54, []byte{17, 165, 134, 138}, 32, "smart rollup commitment hash"},
	{[]byte(PrefixSmartRollupState), 54, []byte{17, 165, 235, 240}, 32, "smart rollup state hash"},
	{[]byte(PrefixSmartRollupInbox), 55, []byte{3, 255, 138, 145, 110}, 32, "smart rollup inbox hash"},
	{[]byte(PrefixSmartRollupMerkelizedTree), 55, []byte{3, 255, 138, 145, 140}, 32, "smart rollup merkelized tree hash"},
}

func getBase58EncodingForDecode(data []byte) (base58Encoding, error) {
	for _, e := range base58Encodings {
		if len(data) != e.DecodedLength+len(e.DecodedPrefix) {
			continue
		}
		found := true
		for i := range e.DecodedPrefix {
			if e.DecodedPrefix[i] != data[i] {
				found = false
				break
			}
		}
		if found {
			return e, nil
		}
	}
	return base58Encoding{}, errors.New("Unknown base58 encoding")
}

func getBase58EncodingForEncode(data, prefix []byte) (base58Encoding, error) {
	for _, e := range base58Encodings {
		if len(data) != e.DecodedLength {
			continue
		}
		found := true
		for i := range prefix {
			if e.EncodedPrefix[i] != prefix[i] {
				found = false
				break
			}
		}
		if found {
			return e, nil
		}
	}
	return base58Encoding{}, errors.New("Unknown base58 encoding")
}

// DecodeBase58 -
func DecodeBase58(data string) ([]byte, error) {
	decoded, err := base58Enc.CheckDecode(data)
	if err != nil {
		return nil, err
	}
	enc, err := getBase58EncodingForDecode(decoded)
	if err != nil {
		return nil, err
	}

	return decoded[len(enc.DecodedPrefix):], nil
}

// DecodeBase58ToString -
func DecodeBase58ToString(data string) (string, error) {
	decoded, err := base58Enc.CheckDecode(data)
	if err != nil {
		return "", err
	}
	enc, err := getBase58EncodingForDecode(decoded)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(decoded[len(enc.DecodedPrefix):]), nil
}

// MustDecodeBase58 -
func MustDecodeBase58(data string) []byte {
	decoded, err := base58Enc.CheckDecode(data)
	if err != nil {
		panic(err)
	}
	enc, err := getBase58EncodingForDecode(decoded)
	if err != nil {
		panic(err)
	}

	return decoded[len(enc.DecodedPrefix):]
}

// DecodeBase58String -
func DecodeBase58String(data string) (string, error) {
	b, err := DecodeBase58(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// EncodeBase58 -
func EncodeBase58(data, prefix []byte) (string, error) {
	enc, err := getBase58EncodingForEncode(data, prefix)
	if err != nil {
		return "", err
	}
	return base58Enc.CheckEncode(append(enc.DecodedPrefix, data...)), nil
}

// EncodeBase58String -
func EncodeBase58String(data string, prefix []byte) (string, error) {
	b, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	return EncodeBase58(b, prefix)
}

// MustEncodeBase58 -
func MustEncodeBase58(data, prefix []byte) string {
	enc, err := getBase58EncodingForEncode(data, prefix)
	if err != nil {
		panic(err)
	}
	return base58Enc.CheckEncode(append(enc.DecodedPrefix, data...))
}

// MustEncodeOperationHash -
func MustEncodeOperationHash(data []byte) string {
	return MustEncodeBase58(data, []byte(PrefixOperationHash))
}
