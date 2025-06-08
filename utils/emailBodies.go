package utils

import "fmt"

func BuildResetEmailBody(domain, token string) string {
	return fmt.Sprintf(`
<html>
  <body style="margin: 0; padding: 0; background-color: #f8f9fa;">
    <table width="100%%" border="0" cellspacing="0" cellpadding="0" style="font-family: Arial, sans-serif;">
      <tr>
        <td align="center" style="padding: 20px 0;">
          <table width="600" border="0" cellspacing="0" cellpadding="0" style="background-color: #e2e8f0; border-radius: 8px; padding: 20px;">
            <tr>
              <td align="center" style="padding: 20px 40px;">
                <h1 style="font-size: 24px; font-weight: bold; color: #1a202c;">
                  Tilbakestill ditt passord
                </h1>
                <p style="font-size: 16px; color: #4a5568; margin-bottom: 20px;">
                  Trykk på knappen under eller lim inn denne URL-en i din nettleser:
                </p>
                <a href="%[2]s/auth/glemt-passord/%[1]s" style="font-size: 16px; color: #3182ce; text-decoration: underline; word-wrap: break-word;">
                  %[2]s/auth/glemt-passord/%[1]s
                </a>
                <div style="margin-top: 20px;">
                  <a href="%[2]s/auth/glemt-passord/%[1]s"
                    style="display: inline-block; padding: 10px 20px; font-size: 16px; color: #fff; background-color: #2d3748; border-radius: 4px; text-decoration: none;">
                    Tilbakestill passord
                  </a>
                </div>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>
`, token, domain)
}
