package field

// Arithmetic Expression Operators
const (
	Abs      = "$abs"
	Add      = "$add"
	Ceil     = "$ceil"
	Divide   = "$divide"
	Exp      = "$exp"
	Floor    = "$floor"
	Ln       = "$ln"
	Log      = "$log"
	Log10    = "$log10"
	Multiply = "$multiply"
	Pow      = "$pow"
	Round    = "$round"
	Sqrt     = "$sqrt"
	Subtract = "$subtract"
	Trunc    = "$trunc"
)

// Array Expression Operators
const (
	ArrayToObject = "$arrayToObject"
	ConcatArrays  = "$concatArrays"
	Filter        = "$filter"
	IndexOfArray  = "$indexOfArray"
	IsArray       = "$isArray"
	Map           = "$map"
	ObjectToArray = "$objectToArray"
	Range         = "$range"
	Reduce        = "$reduce"
	ReverseArray  = "$reverseArray"
	Zip           = "$zip"
)

const (
	Cmp = "$cmp"
)

// Conditional Expression Operators
const (
	Cond   = "$cond"
	IfNull = "$ifNull"
	Switch = "$switch"
)

// Date Expression Operators
const (
	DateFromParts  = "$dateFromParts"
	DateFromString = "$dateFromString"
	DateToParts    = "$dateToParts"
	DateToString   = "$dateToString"
	DayOfMonth     = "$dayOfMonth"
	DayOfWeek      = "$dayOfWeek"
	DayOfYear      = "$dayOfYear"
	Hour           = "$hour"
	IsoDayOfWeek   = "$isoDayOfWeek"
	IsoWeek        = "$isoWeek"
	IsoWeekYear    = "$isoWeekYear"
	Millisecond    = "$millisecond"
	Minute         = "$minute"
	Month          = "$month"
	Second         = "$second"
	ToDate         = "$toDate"
	Week           = "$week"
	Year           = "$year"
)

// Literal Expression Operator
const (
	Literal = "$literal"
)

// Object Expression Operators
const (
	MergeObjects = "$mergeObjects"
)

// Set Expression Operators
const (
	AllElementsTrue = "$allElementsTrue"
	AnyElementTrue  = "$anyElementTrue"
	SetDifference   = "$setDifference"
	SetEquals       = "$setEquals"
	SetIntersection = "$setIntersection"
	SetIsSubset     = "$setIsSubset"
	SetUnion        = "$setUnion"
)

// String Expression Operators
const (
	Concat       = "$concat"
	IndexOfBytes = "$indexOfBytes"
	IndexOfCP    = "$indexOfCP"
	Ltrim        = "$ltrim"
	RegexFind    = "$regexFind"
	RegexFindAll = "$regexFindAll"
	RegexMatch   = "$regexMatch"
	Rtrim        = "$rtrim"
	Split        = "$split"
	StrLenBytes  = "$strLenBytes"
	StrLenCP     = "$strLenCP"
	Strcasecmp   = "$strcasecmp"
	Substr       = "$substr"
	SubstrBytes  = "$substrBytes"
	SubstrCP     = "$substrCP"
	ToLower      = "$toLower"
	ToString     = "$toString"
	Trim         = "$trim"
	ToUpper      = "$toUpper"
)

// Trigonometry Expression Operators
const (
	Sin              = "$sin"
	Cos              = "$cos"
	Tan              = "$tan"
	Asin             = "$asin"
	Acos             = "$acos"
	Atan             = "$atan"
	Atan2            = "$atan2"
	Asinh            = "$asinh"
	Acosh            = "$acosh"
	Atanh            = "$atanh"
	DegreesToRadians = "$degreesToRadians"
	RadiansToDegrees = "$radiansToDegrees"
)

// Type Expression Operators
const (
	Convert    = "$convert"
	ToBool     = "$toBool"
	ToDecimal  = "$toDecimal"
	ToDouble   = "$toDouble"
	ToInt      = "$toInt"
	ToLong     = "$toLong"
	ToObjectID = "$toObjectId"
)

// Accumulators ($group)
const (
	Avg        = "$avg"
	First      = "$first"
	Last       = "$last"
	StdDevPop  = "$stdDevPop"
	StdDevSamp = "$stdDevSamp"
	Sum        = "$sum"
)

// Variable Expression Operators
const (
	Let = "$let"
)
