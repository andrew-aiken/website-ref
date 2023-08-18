import boto3
import logging
import os

REGION = os.environ['AWS_REGION']

EC2 = boto3.client('ec2', region_name=REGION)
ASG = boto3.client('autoscaling', region_name=REGION)

log_level = os.environ.get('LOGLEVEL', 'INFO').upper()

logging.basicConfig(
    level=log_level,
    format='%(levelname)s   %(asctime)s   %(message)s'
)

logger = logging.getLogger()
logger.setLevel(log_level)

def get_asg_name(instance_id):
    try:
        response = EC2.describe_instances(
            InstanceIds = [instance_id]
        )
        for tag_pair in response["Reservations"][0]["Instances"][0]["Tags"]:
            if tag_pair["Key"] == "aws:autoscaling:groupName":
                logger.debug(f'asg name {tag_pair["Value"]}')
                return tag_pair["Value"]
    except Exception as error:
        logger.critical(f'Error when trying to get asg from {instance_id}')
        logger.error(error)
    return None

def get_lifecycle_hook_name(asg_name):
    try:
        response = ASG.describe_lifecycle_hooks(
            AutoScalingGroupName = asg_name
        )
        logger.debug(f'lifecycle hook name {response["LifecycleHooks"][0]["LifecycleHookName"]}')
        return response["LifecycleHooks"][0]["LifecycleHookName"]
    except Exception as error:
        logger.critical(f'Error when trying to get lifecycle hook associated with {asg_name}')
        logger.error(error)
    return None

def lambda_handler(event, context):
    logger.debug(event["detail"])

    if event["detail"]["status"] != "DRAINING":
        logger.critical('Status is not "DRAINING"')
        return
    elif event["detail"]["pendingTasksCount"] != 0:
        logger.critical('There are pending tasks on the node')
        return
    elif event["detail"]["runningTasksCount"] != 0:
        logger.critical('There are still running tasks on the node')
        return

    asg_name = get_asg_name(instance_id=event["detail"]["ec2InstanceId"])

    if asg_name != None:
        life_cycle_hook_name = get_lifecycle_hook_name(asg_name=asg_name)

    if life_cycle_hook_name != None:
        try:
            ASG.complete_lifecycle_action(
                LifecycleHookName=life_cycle_hook_name,
                AutoScalingGroupName=asg_name,
                LifecycleActionResult='CONTINUE',
                InstanceId=event["detail"]["ec2InstanceId"]
            )
            logger.info(f'Marking lifecycle to remove {event["detail"]["ec2InstanceId"]}')
        except Exception as error:
            logger.critical(f'Error removing {event["detail"]["ec2InstanceId"]}')
            logger.error(error)
    return
