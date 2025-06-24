// Luhn algorithm
package luhn

func IsValidLuhn(number string) bool {
	if number == "" {
		return false
	}

	sum := 0
	isOddLength := len(number)%2 != 0

	for i, r := range number {
		if r < '0' || r > '9' {
			return false
		}

		n := int(r - '0')
		if (i%2 == 0) == !isOddLength { // чет-нечет сдвиг
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
	}

	return sum%10 == 0
}
