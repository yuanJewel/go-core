package email

import (
	"html/template"
	"strings"
)

func Text2HTML(data string) string {
	escapedData := template.HTMLEscapeString(data)
	replacer := strings.NewReplacer(
		"\n", "<br>",
		" ", "&nbsp;",
		"\t", "&nbsp;&nbsp;&nbsp;&nbsp;",
		"\r", "",
	)

	return replacer.Replace(escapedData)
}

func Table2HTML(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	var tableBuilder strings.Builder
	tableBuilder.WriteString(`<table border="1" cellpadding="5" cellspacing="0" style="border-collapse: collapse; width: 100%; max-width: 100%;">`)

	// 生成表头
	tableBuilder.WriteString("<thead><tr>")
	for _, header := range rows[0] {
		tableBuilder.WriteString(`<th style="background-color: #f2f2f2; padding: 8px; text-align: left; word-break: break-word; white-space: nowrap;">`)
		tableBuilder.WriteString(template.HTMLEscapeString(header))
		tableBuilder.WriteString("</th>")
	}
	tableBuilder.WriteString("</tr></thead>")

	// 生成表格主体
	tableBuilder.WriteString("<tbody>")
	for i := 1; i < len(rows); i++ {
		tableBuilder.WriteString("<tr>")
		for _, cell := range rows[i] {
			tableBuilder.WriteString(`<td style="padding: 8px; word-wrap: break-word; overflow-wrap: break-word; word-break: break-word;">`)

			// 处理单元格中的换行符
			escapedCell := template.HTMLEscapeString(cell)
			lines := strings.Split(escapedCell, "\n")

			if len(lines) == 1 {
				tableBuilder.WriteString(lines[0])
			} else {
				for j, line := range lines {
					tableBuilder.WriteString(`<div style="margin-top: 2px; padding-left: 12px; position: relative;">`)
					if j == 0 {
						tableBuilder.WriteString(`<span style="color: #222; position: absolute; left: 0; font-size: 10px;">→</span>`)
					} else {
						tableBuilder.WriteString(`<span style="color: #222; position: absolute; left: 0; font-size: 10px;">⤷</span>`)
					}
					tableBuilder.WriteString(line)
					tableBuilder.WriteString("</div>")
				}
			}

			tableBuilder.WriteString("</td>")
		}
		tableBuilder.WriteString("</tr>")
	}
	tableBuilder.WriteString("</tbody></table>")

	return tableBuilder.String()
}
