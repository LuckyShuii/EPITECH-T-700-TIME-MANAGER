package template

import "fmt"

func BaseMailTemplate(title, content, buttonLabel, buttonLink string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		/* Hover: black background, white text */
		.btn:hover { background:#000000 !important; color:#ffffff !important; }
	</style>
</head>

<body style="margin:0;padding:0;background:#ffffff;font-family:Arial, Helvetica, sans-serif;color:#000;">
	<table width="100%%" cellpadding="0" cellspacing="0" style="background:#ffffff;">
		<tr>
			<td align="center" style="padding:24px 12px;">

				<table width="600" cellpadding="0" cellspacing="0"
					style="width:600px;max-width:600px;background:#ffffff;border:3px solid #000000;padding:24px;">

					<tr>
						<td align="center"
							style="font-size:22px;font-weight:700;letter-spacing:-0.2px;
							line-height:1.2;padding-bottom:16px;">
							%s
						</td>
					</tr>

					<tr>
						<td style="font-size:14px;line-height:1.7;padding-bottom:20px;">
							%s
						</td>
					</tr>

					%s

					<tr>
						<td style="padding-top:20px;font-size:12px;line-height:1.6;text-align:center;">
							<div style="font-weight:700;">TimeManager</div>
							<div>Automatic email • Please do not reply</div>
						</td>
					</tr>

				</table>

			</td>
		</tr>
	</table>
</body>
</html>
`,
		title,
		title,
		content,
		func() string {
			if buttonLabel == "" {
				return ""
			}

			return fmt.Sprintf(`
<tr>
	<td align="center" style="padding-bottom:8px;">
		<table cellpadding="0" cellspacing="0" border="0">
			<tr>
				<td align="center">
					<a class="btn" href="%s"
						style="display:inline-block;padding:12px 16px;
						border:3px solid #000000;background:#ffffff;color:#000000;
						text-decoration:none;font-weight:700;font-size:14px;line-height:1;
						text-transform:uppercase;letter-spacing:0.6px;">
						%s
					</a>
				</td>
			</tr>
		</table>
	</td>
</tr>
`, buttonLink, buttonLabel)
		}(),
	)
}
