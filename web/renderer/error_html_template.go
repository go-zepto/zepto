package renderer

var ERROR_HTML_TEMPLATE = `
<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Internal Error - {{ .Title }}</title>
  		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			*,
			*::before,
			*::after {
				box-sizing: border-box;
			}
			body {
				margin: 0;
				padding: 0;
				font-family: Arial, serif;
			}
			table {
				border-collapse: collapse;
				width: 100%;
				overflow-wrap: break-word;
			}
			.table-key {
				width: 240px;
			}
			@media only screen and (max-width: 600px) {
			  	table {
					overflow: auto;
				}
				.table-key {
					width: auto;
				}
			}
			table {
				border: 1px;
				color: #666;
			}
			table, th, td {
				border: 1px solid #c6c6c6;
			}
			th, td {
				padding: 16px;
			}
			table {
				margin-bottom: 32px;
				font-size: 14px;
				table-layout: fixed;
			}
			.container {
				padding: 0 32px;
			}
			.topbar {
				height: 48px;
				width: 100%;
				background-color: #0984E3;
				padding: 0 32px;
			}
			.topbar img {
				height: 16px;
				margin-top: 16px;
			}
			.topbar .links {
				float: right;
				margin-top: 16px;
				color: #FFF !important;
				text-decoration: none;
			}
			.topbar .links a {
				color: #FFF;
				text-decoration: none;
				margin-left: 8px;
			}
			.error-header {
				width: 100%;
				background-color: #f7f7f7;
				padding: 48px 32px;
				color: #0984E3;
				vertical-align: middle;
				border-bottom: 1px solid #e9e9e9;
			}
			.error-header .error-title {
				margin-top: 0;
				font-size: 32px;
			}
			.error-header .error-info {
				margin-bottom: 8px;
			}
			.error-header .error-info span {
				color: #666666;
			}
			.error-details {
				padding-top: 32px;
			}
			.table-title {
				color: #0984E3;
				margin-bottom: 8px;
			}
			.error-stack-box {
				padding: 24px;
				border: 1px solid #c6c6c6;
				line-height: 24px;
				max-height: 300px;
				overflow: auto;
				color: #666;
			}
		</style>
	</head>

	<body>
		<div class="topbar">
			<img alt="Zepto" src="data:image/png;base64, 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" />
			<div class="links">
				<a href="https://go-zepto.github.io/zepto" target="_blank">
					Documentation
				</a>
				<a href="https://godoc.org/github.com/go-zepto/zepto" target="_blank">
					GoDoc
				</a>
				<a href="https://github.com/go-zepto/zepto" target="_blank">
					Github
				</a>
			</div>
		</div>
		<div class="error-header">
			<div class="error-info">
				Error <span>at {{ .Req.Method }} {{ .Req.URL.Path }}</span>
			</div>
			<div class="error-title">
				{{ .Title }}
			</div>
		</div>

		<div class="container error-details">
			<div class="table-title">Error Trace</div>
			<div class="error-stack-box">
				{{ .Trace }}
			</div>
		</div>
		<div class="container error-details">
			<div class="table-title">Request Info</div>
			<table>
				<tr>
					<td class="table-key">
						<strong>URI</strong>
					</td>
					<td>
						{{ .URI }}
					</td>
				</tr>
				<tr>
					<td class="table-key">
						<strong>Query String</strong>
					</td>
					<td>
						{{ .Req.URL.RawQuery }}
					</td>
				</tr>
			</table>
			<div class="table-title">Headers</div>
			<table>
				{{ range $key, $value := .Req.Header }}
					<tr>
						<td class="table-key">
							<strong>{{ $key }}</strong>
						</td>
						<td>
							{{ $value }}
						</td>
					</tr>
				{{ end }}
			</table>
		</div>
	</body>
</html>
`
