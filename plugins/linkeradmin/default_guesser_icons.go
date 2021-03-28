package linkeradmin

// Map of Material UI Icons -> Possible Resource Names
//
// Official Icon Repository: https://material.io/resources/icons
var defaultGuesserIcons = map[string][]string{
	"account_balance_wallet": {
		"wallet",
	},
	"alarm": {
		"timer",
		"clock",
	},
	"all_inbox": {
		"inbox",
	},
	"announcement": {
		"info",
	},
	"article": {
		"news",
	},
	"assignment": {},
	"backup":     {},
	"book": {
		"publication",
		"volume",
		"title",
		"novel",
		"class",
		"course",
	},
	"bug_report": {
		"bug",
	},
	"build": {},
	"card_giftcard": {
		"gift",
		"giftcard",
	},
	"check_circle": {
		"check",
		"approval",
	},
	"circle_notifications": {
		"notification",
	},
	"code": {
		"sourcecode",
	},
	"credit_card": {
		"payment_method",
		"payment",
		"subscription",
	},
	"dashboard": {},
	"event": {
		"agenda",
		"meet",
		"calendar",
	},
	"extension": {
		"plugin",
		"addon",
		"module",
	},
	"favorite": {
		"save",
	},
	"file_present": {},
	"gif":          {},
	"help":         {},
	"history":      {},
	"home": {
		"house",
	},
	"http": {},
	"lock": {},
	"info": {
		"information",
	},
	"label": {
		"tag",
	},
	"language": {
		"world",
		"country",
	},
	"nightlight_round": {
		"nigh",
	},
	"paid": {
		"payment",
	},
	"pending": {
		"pending",
	},
	"perm_contact": {
		"contact",
	},
	"perm_media": {
		"media",
	},
	"perm_phone_msg": {
		"phonemessage",
		"phonemsg",
	},
	"pet": {
		"dog",
		"cat",
	},
	"preview": {},
	"print":   {},
	"question_answer": {
		"conversation",
		"topic",
	},
	"schedule": {},
	"search":   {},
	"segment":  {},
	"settings": {
		"setting",
		"configuration",
		"config",
		"preference",
	},
	"shopping_card": {
		"shopping_cart",
		"cart",
	},
	"speaker_notes": {},
	"star_rate": {
		"rate",
	},
	"sticky_note_2": {
		"note",
	},
	"store": {
		"ecommerce",
	},
	"support":  {},
	"verified": {},
	"work": {
		"job",
	},
	"error":   {},
	"warning": {},
	"album":   {},
	"games": {
		"game",
	},
	"movie": {
		"cine",
		"cinema",
	},
	"playlist_play": {
		"playlist",
	},
	"videocam": {
		"live",
		"webinar",
		"record",
	},
	"web": {
		"landingpage",
		"website",
	},
	"comment": {},
	"call":    {},
	"chat": {
		"chatmessage",
	},
	"domain_verified": {},
	"email": {
		"mail",
	},
	"forum":     {},
	"message":   {},
	"qr_code":   {},
	"sms":       {},
	"archive":   {},
	"block":     {},
	"draft":     {},
	"flag":      {},
	"inbox":     {},
	"inventory": {},
	"link": {
		"url",
	},
	"push_pin": {
		"pin",
	},
	"airplane_ticket": {
		"travel",
		"flight",
	},
	"cable":        {},
	"credit_score": {},
	"dark_mode": {
		"night",
	},
	"device_thermostat": {
		"thermostat",
	},
	"devices": {
		"device",
	},
	"lightmode": {
		"weather",
	},
	"pattern":      {},
	"price_change": {},
	"quiz":         {},
	"review":       {},
	"task":         {},
	"widgets":      {},
	"attach_money": {
		"currency",
	},
	"checklist": {},
	"insert_chart": {
		"chart",
		"statistics",
		"analytics",
	},
	"insert_photo": {
		"image",
		"photo",
		"figure",
	},
	"attachment": {},
	"cloud":      {},
	"cloud_download": {
		"download",
	},
	"cloud_upload": {
		"upload",
	},
	"folder": {
		"directory",
	},
	"computer": {
		"desktop",
	},
	"memory":    {},
	"monitor":   {},
	"headphone": {},
	"headset":   {},
	"keyboard":  {},
	"laptop":    {},
	"smartphone": {
		"phone",
	},
	"toy": {},
	"tv": {
		"television",
	},
	"videogame": {},
	"watch":     {},
	"assistant": {},
	"audiotrack": {
		"music",
	},
	"cases": {
		"case",
	},
	"pallete": {
		"color",
	},
	"camera":      {},
	"slideshow":   {},
	"style":       {},
	"attractions": {},
	"badge":       {},
	"car_rental":  {},
	"car_repair":  {},
	"celebration": {},
	"cleaning_services": {
		"clean",
		"cleanning",
	},
	"delivery_dining": {
		"delivery",
	},
	"direction": {},
	"directions_bike": {
		"bicycle",
		"bike",
		"cyclist",
		"rider",
	},
	"directions_boat": {
		"boat",
		"ship",
	},
	"directions_bus": {
		"bus",
	},
	"directions_car": {
		"car",
	},
	"directions_subway": {
		"subway",
		"train",
	},
	"directions_run": {
		"walk",
		"run",
	},
	"dry_cleaning": {},
	"fastfood": {
		"food",
	},
	"festival": {},
	"home_repair_service": {
		"repair",
	},
	"hotel": {
		"hostel",
		"motel",
	},
	"layers": {
		"layer",
	},
	"local_activity": {
		"ticket",
	},
	"local_bar": {
		"drink",
	},
	"local_cafe": {
		"cafe",
		"coffe",
	},
	"local_dinning": {
		"dinning",
	},
	"local_hospital": {
		"hospital",
		"health",
	},
	"local_gas_station": {
		"gas",
		"gasstation",
	},
	"local_grocery_store": {
		"market",
		"supermarket",
	},
	"local_laundry_service": {
		"laundry",
	},
	"local_library": {
		"library",
	},
	"local_mall": {
		"mall",
		"shoppingmall",
	},
	"local_offer": {
		"offer",
	},
	"local_pharmacy": {
		"pharmacy",
	},
	"local_pizza": {
		"pizza",
	},
	"local_police": {
		"police",
	},
	"local_shipping": {
		"shipping",
	},
	"local_taxi": {
		"taxi",
	},
	"map": {
		"guide",
	},
	"museum":     {},
	"navigation": {},
	"my_location": {
		"mylocation",
		"currentlocation",
		"target",
	},
	"park": {
		"florest",
	},
	"place": {
		"adress",
		"site",
		"locale",
		"spot",
		"position",
	},
	"terrain": {
		"mountain",
	},
	"theater_comedy": {
		"comedy",
	},
	"traffic": {},
	"two_wheeler": {
		"motorcycle",
	},
	"wine_bar": {
		"wine",
	},
	"child_care": {
		"child",
	},
	"child_friendly": {
		"baby",
	},
	"fitness_center": {
		"gym",
	},
	"kitchen": {
		"kitchen",
	},
	"bathroom": {},
	"bed":      {},
	"feed":     {},
	"cake": {
		"birthday",
	},
	"construction": {},
	"coronavirus": {
		"virus",
	},
	"emoji_emotions": {
		"emoji",
	},
	"emoji_events": {
		"championship",
	},
	"emoji_objects": {
		"idea",
	},
	"emoji_symbols": {
		"symbol",
	},
	"engineering": {
		"engineer",
	},
	"group": {
		"person",
		"user",
		"author",
		"admin",
		"employee",
	},
	"masks": {
		"mask",
	},
	"luggage":    {},
	"poll":       {},
	"psychology": {},
	"public":     {},
	"recommend": {
		"like",
	},
	"school": {
		"university",
		"college",
	},
	"science": {},
	"share":   {},
}

// Generate a map resourceName -> Material UI Icon
func generateInverseDefaultGuesserIcons(guesserIcons map[string][]string) map[string]string {
	inverseMap := make(map[string]string)
	for materialIconName, possibleResources := range guesserIcons {
		inverseMap[materialIconName] = materialIconName
		for _, pr := range possibleResources {
			inverseMap[pr] = materialIconName
		}
	}
	return inverseMap
}
