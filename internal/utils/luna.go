package utils

const (
	digitThreshold = 9
)

// IsValidLuna валидация номера заказа по алгоритму Луна
func IsValidLuna(number string) bool {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if digit < 0 || digit > 9 {
			return false
		}

		if alternate {
			digit *= 2
			if digit > digitThreshold {
				digit = (digit % 10) + 1
			}
		}

		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}
