package validations

import (
	"reflect"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/go-playground/validator/v10"
	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/helpers"
)

const (
	defaultMaxSize = 10
)

// Register -
func Register(v *validator.Validate, cfg config.APIConfig) error {
	if err := v.RegisterValidation("address", addressValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("contract", contractValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("opg", opgValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("network", networkValidator(cfg.Networks)); err != nil {
		return err
	}

	if err := v.RegisterValidation("status", statusValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("faversion", faVersionValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("fill_type", fillTypeValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("search", searchStringValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("gt_int64_ptr", greatThanInt64PtrValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("bcd_max_size", maxSizeValidator(cfg.PageSize)); err != nil {
		return err
	}

	if err := v.RegisterValidation("global_constant", constantAddressValidator()); err != nil {
		return err
	}

	if err := v.RegisterValidation("smart_rollup", smartRollupValidator()); err != nil {
		return err
	}

	return nil
}

func addressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return bcd.IsAddress(fl.Field().String())
	}
}

func smartRollupValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return bcd.IsSmartRollupHash(fl.Field().String())
	}
}

func contractValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return bcd.IsContract(fl.Field().String())
	}
}

func networkValidator(networks []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		network := fl.Field().String()
		return helpers.StringInArray(network, networks)
	}
}

func opgValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		hash := fl.Field().String()
		if !strings.HasPrefix(hash, "o") || len(hash) != 51 {
			return false
		}
		_, _, err := base58.CheckDecode(hash)
		return err == nil
	}
}

func statusValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		status := fl.Field().String()
		data := strings.Split(status, ",")
		for i := range data {
			if !helpers.StringInArray(data[i], []string{
				consts.Applied,
				consts.Backtracked,
				consts.Failed,
				consts.Skipped,
			}) {
				return false
			}
		}
		return true
	}
}

func faVersionValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		version := fl.Field().String()
		return helpers.StringInArray(version, []string{
			consts.FA1Tag,
			"fa12",
			consts.FA2Tag,
		})
	}
}

func fillTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		fillType := fl.Field().String()
		return helpers.StringInArray(fillType, []string{
			"empty",
			"current",
			"initial",
		})
	}
}

func searchStringValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 2 && len(fl.Field().String()) < 256
	}
}

func maxSizeValidator(maxSize uint64) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if maxSize == 0 {
			maxSize = defaultMaxSize
		}
		fieldValue := fl.Field().Int()
		if fieldValue < 0 {
			return false
		}
		return uint64(fieldValue) <= maxSize
	}
}

func greatThanInt64PtrValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()
		kind := field.Kind()

		currentField, currentKind, _, ok := fl.GetStructFieldOK2()
		if !ok {
			return false
		}

		switch {
		case kind == reflect.Ptr && currentKind == reflect.Ptr:
			return true
		case kind == reflect.Ptr && currentKind == reflect.Int64:
			return true
		case kind == reflect.Int64 && currentKind == reflect.Ptr:
			return true
		case kind == reflect.Int64 && currentKind == reflect.Int64:
			return field.Int() > currentField.Int()
		default:
			return false
		}
	}
}

func constantAddressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		address := fl.Field().String()
		if !strings.HasPrefix(address, "expr") || len(address) != 54 {
			return false
		}
		_, _, err := base58.CheckDecode(address)
		return err == nil
	}
}
