package cmd

const defAppFieldValue = "undefined"

var (
	version = defAppFieldValue
	builded = defAppFieldValue
	hash    = defAppFieldValue

	AppInfo = version + ", builded: " + builded + ", hash: " + hash
)
