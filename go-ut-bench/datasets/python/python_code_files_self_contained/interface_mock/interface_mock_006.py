import json
import smtplib
SMTP_SERVER = 'smtp.gmail.com'
SMTP_PORT = 587
EMAIL_ADDRESS = 'your.email@gmail.com'
EMAIL_PASSWORD = 'your.password'

def task_func(input_data=None, smtp_server=SMTP_SERVER, smtp_port=SMTP_PORT, email_address=EMAIL_ADDRESS, email_password=EMAIL_PASSWORD, smtp=None):
    if input_data is None:
        return []
    try:
        data = json.loads(input_data)
        recipient_email = data.get('recipient')
        names = data.get('names', [])
    except (json.JSONDecodeError, ValueError):
        return []
    if not recipient_email or not names:
        return []
    message = 'Subject: Extracted Names\n\n' + '\n'.join(names)
    if smtp:
        server = smtp(smtp_server, smtp_port)
    else:
        server = smtplib.SMTP(smtp_server, smtp_port)
    server.starttls()
    server.login(email_address, email_password)
    server.sendmail(email_address, recipient_email, message)
    server.quit()
    return names
