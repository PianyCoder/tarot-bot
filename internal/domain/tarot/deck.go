package tarot

type Deck map[int]Card

func NewDeck() Deck {
	deck := make(Deck)

	// Старшие Арканы
	deck[0] = Card{ID: 0, Title: "I. Маг", Description: "контроль, инициатива, сила влияния"}
	deck[1] = Card{ID: 1, Title: "II. Жрица", Description: "интуиция, тайна, скрытая женская сила"}
	deck[2] = Card{ID: 2, Title: "III. Императрица", Description: "женская природа, чувственность, забота"}
	deck[3] = Card{ID: 3, Title: "IV. Император", Description: "контроль, структура, власть, правила"}
	deck[4] = Card{ID: 4, Title: "V. Верховный Жрец (Иерофант)", Description: "традиции, долг, мораль, воспитание"}
	deck[5] = Card{ID: 5, Title: "VI. Влюблённые", Description: "выбор, страсть, зависимость"}
	deck[6] = Card{ID: 6, Title: "VII. Колесница", Description: "стремление, движение, сила воли"}
	deck[7] = Card{ID: 7, Title: "VIII. Справедливость (Правосудие)", Description: "равновесие, расплата, последствия"}
	deck[8] = Card{ID: 8, Title: "IX. Отшельник", Description: "одиночество, самопознание, дистанция"}
	deck[9] = Card{ID: 9, Title: "X. Колесо Фортуны", Description: "перемены, случайности, циклы"}
	deck[10] = Card{ID: 10, Title: "XI. Сила", Description: "страсть, внутренняя мощь, обуздание"}
	deck[11] = Card{ID: 11, Title: "XII. Повешенный", Description: "жертва, зависание, переосмысление"}
	deck[12] = Card{ID: 12, Title: "XIII. Смерть", Description: "конец, трансформация, точка"}
	deck[13] = Card{ID: 13, Title: "XIV. Умеренность", Description: "гармония, терпение, принятие"}
	deck[14] = Card{ID: 14, Title: "XV. Дьявол", Description: "зависимость, страсть, одержимость"}
	deck[15] = Card{ID: 15, Title: "XVI. Башня", Description: "разрушение, шок, срыв"}
	deck[16] = Card{ID: 16, Title: "XVII. Звезда", Description: "надежда, красота, вдохновение"}
	deck[17] = Card{ID: 17, Title: "XVIII. Луна", Description: "иллюзии, страхи, скрытое"}
	deck[18] = Card{ID: 18, Title: "XIX. Солнце", Description: "радость, открытость, признание"}
	deck[19] = Card{ID: 19, Title: "XX. Суд (Страшный суд)", Description: "возрождение, карма, поворот"}
	deck[20] = Card{ID: 20, Title: "XXI. Мир", Description: "завершение, цель, целостность"}
	deck[21] = Card{ID: 21, Title: "0. Шут", Description: "свобода, импульс, наивность, начало"}

	return deck
}
