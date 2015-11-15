# Mondo Hackathon Email Receiver

This work is done for the Mondo Hackathon.

This is a small service that connects to your Gmail account and fetches a
specific email, givin in the URL.

```
localhost:3000/emails/<ID>
```

The ID can be obtained in the Gmail web UI by taking the last digits from the
URL (`15108647a6a22d93`).

## Attachments

Images can be sent over as attachments. This means that some images won't be
linked properly in the email, but they will link to an attachment.

We handle this by looking if there are available attachments and if so, replace
them by their base64 value in the HTML output.

## Webhook

You can configure a webhook to receive the email data (see the email package).
