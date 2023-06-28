package list

const HtmlHead = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>Mod List</title>
<style>
body {
	font-family: sans-serif;
}
</style>
</head>
<body>
<h1>Mod List</h1>
<ul>`

const HtmlFoot = `
</ul>
</body>
</html>`

func Html() (string, error) {
	urls, err := CreateModList()

	if err != nil {
		return "", err
	}

	out := HtmlHead

	for _, url := range urls {
		out += "<li><a href=\"" + url.Url + "\">" + url.Mod + "</a></li>"
	}

	out += HtmlFoot

	return out, nil
}
