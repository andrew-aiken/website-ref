# UML Cloud Computing Club | AWS Event-driven

Most of this repo is just example Terraform code. The only real file of note is the presentation pdf.


## Samples

Trigger a sqs queue from the cli

```bash
aws sqs send-message \
    --queue-url https://sqs.us-east-1.amazonaws.com/683454754281/example-queue \
    --message-body '{"message": "Trigger from SQS"}'
```
