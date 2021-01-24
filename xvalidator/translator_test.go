package xvalidator

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
	"unsafe"
)

func TestTranslator(t *testing.T) {
	v := validator.New()
	type testStruct struct {
		String string `validate:"required"`
	}

	for _, tc := range []struct {
		giveTranslator   locales.Translator
		giveRegisterFn   TranslationRegisterHandler
		wantRequiredText string
	}{
		{EnLocaleTranslator(), EnTranslationRegisterFunc(), "String is a required field"},
		{FrLocaleTranslator(), FrTranslationRegisterFunc(), "String est un champ obligatoire"},
		{JaLocaleTranslator(), JaTranslationRegisterFunc(), "Stringは必須フィールドです"},
		{ZhLocaleTranslator(), ZhTranslationRegisterFunc(), "String为必填字段"},
		{ZhHantLocaleTranslator(), ZhTwTranslationRegisterFunc(), "String為必填欄位"},
	} {
		translator, err := ApplyTranslationToValidator(v, tc.giveTranslator, tc.giveRegisterFn)
		xtesting.Nil(t, err)
		err = v.Struct(&testStruct{})
		xtesting.NotNil(t, err)
		xtesting.Equal(t, err.(validator.ValidationErrors).Translate(translator)["testStruct.String"], tc.wantRequiredText)
	}

	xtesting.Panic(t, func() { _, _ = ApplyTranslationToValidator(nil, EnLocaleTranslator(), EnTranslationRegisterFunc()) })
	xtesting.Panic(t, func() { _, _ = ApplyTranslationToValidator(v, nil, EnTranslationRegisterFunc()) })
	xtesting.Panic(t, func() { _, _ = ApplyTranslationToValidator(v, EnLocaleTranslator(), nil) })
	_, err := ApplyTranslationToValidator(v, EnLocaleTranslator(), func(v *validator.Validate, trans ut.Translator) error {
		return errors.New("test error")
	})
	xtesting.NotNil(t, err)
	xtesting.Equal(t, err.Error(), "test error")
}

func TestTranslatorRegister(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("test_eq", EqualValidator(0))
	type testStruct struct {
		String string `validate:"required"`
		Int    int    `validate:"test_eq"`
	}

	trans, _ := ApplyTranslationToValidator(val, EnLocaleTranslator(), EnTranslationRegisterFunc())
	fn := RegisterTranslationFunc("required", "required {0}!!!", true)
	_ = val.RegisterTranslation("required", trans, fn, DefaultTranslationFunc())
	_ = val.RegisterTranslation("test_eq", trans, fn, DefaultTranslationFunc()) // <<< error

	err := val.Struct(&testStruct{}).(validator.ValidationErrors)
	xtesting.NotNil(t, err)
	transResults := err.Translate(trans)
	xtesting.Equal(t, transResults["testStruct.String"], "required String!!!")

	err = val.Struct(&testStruct{String: "test", Int: 1}).(validator.ValidationErrors)
	xtesting.NotNil(t, err)
	transResults = err.Translate(trans)
	xtesting.Equal(t, transResults["testStruct.Int"], "Key: 'testStruct.Int' Error:Field validation for 'Int' failed on the 'test_eq' tag")
}

func TestLocaleTranslators(t *testing.T) {
	for _, tc := range []struct {
		give     locales.Translator
		wantName string
	}{
		{EnLocaleTranslator(), "en"},
		{FrLocaleTranslator(), "fr"},
		{IdLocaleTranslator(), "id"},
		{JaLocaleTranslator(), "ja"},
		{NlLocaleTranslator(), "nl"},
		{PtBrLocaleTranslator(), "pt_BR"},
		{RuLocaleTranslator(), "ru"},
		{TrLocaleTranslator(), "tr"},
		{ZhLocaleTranslator(), "zh"},
		{ZhHantLocaleTranslator(), "zh_Hant"},
	} {
		xtesting.Equal(t, tc.give.Locale(), tc.wantName)
	}
}

func TestTranslationRegisterFuncs(t *testing.T) {
	type transText struct {
		text    string
		indexes []int
	}

	for _, tc := range []struct {
		giveFn           TranslationRegisterHandler
		wantFields       bool
		wantRequiredText string
	}{
		{DefaultTranslationRegisterFunc(), false, ""},
		{EnTranslationRegisterFunc(), true, "{0} is a required field"},
		{FrTranslationRegisterFunc(), true, "{0} est un champ obligatoire"},
		{IdTranslationRegisterFunc(), true, "{0} wajib diisi"},
		{JaTranslationRegisterFunc(), true, "{0}は必須フィールドです"},
		{NlTranslationRegisterFunc(), true, "{0} is een verplicht veld"},
		{PtBrTranslationRegisterFunc(), true, "{0} é um campo requerido"},
		{RuTranslationRegisterFunc(), true, "{0} обязательное поле"},
		{TrTranslationRegisterFunc(), true, "{0} zorunlu bir alandır"},
		{ZhTranslationRegisterFunc(), true, "{0}为必填字段"},
		{ZhTwTranslationRegisterFunc(), true, "{0}為必填欄位"},
	} {
		val := validator.New()
		uniTrans := ut.New(EnLocaleTranslator(), EnLocaleTranslator())
		trans, _ := uniTrans.GetTranslator(EnLocaleTranslator().Locale())
		err := tc.giveFn(val, trans)
		xtesting.Nil(t, err)

		field := reflect.ValueOf(trans).Elem().FieldByName("translations")
		fieldValue := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		if tc.wantFields {
			ptr := fieldValue.MapIndex(reflect.ValueOf("required"))
			xtesting.Equal(t, (*transText)(unsafe.Pointer(ptr.Elem().UnsafeAddr())).text, tc.wantRequiredText)
		} else {
			xtesting.Equal(t, fieldValue.Len(), 0)
		}
	}
}
