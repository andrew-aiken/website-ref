import boto3
import json
import logging
import os

CLUSTER = os.environ['ClusterName']
REGION = os.environ['AWS_REGION']

ECS = boto3.client('ecs', region_name=REGION)

log_level = os.environ.get('LOGLEVEL', 'INFO').upper()

logging.basicConfig(
    level=logging.INFO,
    format='%(levelname)s   %(asctime)s   %(message)s'
)

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def find_ecs_instance_info(instance_id):
    try:
        container_instances = ECS.list_container_instances(
            cluster = CLUSTER
        )

        container_instance_descriptions = ECS.describe_container_instances(
            cluster = CLUSTER,
            containerInstances = container_instances['containerInstanceArns']
        )
        for instance in container_instance_descriptions['containerInstances']:
            logger.debug(instance["ec2InstanceId"])
            if instance["ec2InstanceId"] == instance_id:
                return instance["containerInstanceArn"], instance["status"]
    except Exception as error:
        logger.critical('Error when trying to retrieve instance information')
        logger.error(error)
    return None, None

def instance_has_running_tasks(instance_id):
    instance_arn, instance_status = find_ecs_instance_info(instance_id)
    if instance_arn is None:
        logger.critical(f'Could not find instance ID {instance_id}. Letting autoscaling kill the instance.')
        return False
    if instance_status != 'DRAINING':
        try:
            ECS.tag_resource(
                resourceArn = instance_arn,
                tags=[{
                    'key': 'lifecycle-drain',
                    'value': 'true'
                }]
            )
            logger.debug(f'Tagged {instance_arn}')
        except Exception as error:
            logger.critical(f'Could not tag instance {instance_arn}')
            logger.error(error)
        try:
            ECS.update_container_instances_state(
                cluster = CLUSTER,
                containerInstances = [instance_arn],
                status = 'DRAINING'
            )
            logger.info(f'Setting container instance {instance_id} ({instance_arn}) to DRAINING')
        except Exception as error:
            logger.critical(f'Could not set ecs instance ({instance_arn}) to drain')
            logger.error(error)

def lambda_handler(event, context):
    msg = json.loads(event['Records'][0]['Sns']['Message'])
    logger.info(msg)
    if ('LifecycleTransition' not in msg.keys() or \
       msg['LifecycleTransition'].find('autoscaling:EC2_INSTANCE_TERMINATING') == -1):
        logger.critical('Exiting since the lifecycle transition is not EC2_INSTANCE_TERMINATING')
        return
    instance_has_running_tasks(msg['EC2InstanceId'])
