import os

LAMBDA_ARN = os.environ["LAMBDA_ARN"]
PREFIX = "Mackerel::"

def handle_template(request_id, template):
    for _, resource in template.get("Resources", {}).items():
        if resource["Type"].startswith("PREFIX"):
            properties = resource.get("Properties", {})
            properties["ServiceToken"] = LAMBDA_ARN
            resource.update({
                "Type": "Custom::" + resource["Type"][len(PREFIX):],
                "Version": "1.0",
                "Properties": properties,
            })
    return template

def handler(event, context):
    fragment = event["fragment"]
    status = "success"

    try:
        fragment = handle_template(event["requestId"], event["fragment"])
    except Exception:
        status = "failure"

    return {
        "requestId": event["requestId"],
        "status": status,
        "fragment": fragment,
    }
