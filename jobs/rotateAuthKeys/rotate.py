import os
import base64
import datetime
import uuid
from kubernetes import client, config
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa

NAMESPACE = os.environ["NAMESPACE"]
SECRET_NAME = os.environ["SECRET_NAME"]
DEPLOYMENT_NAME = os.environ["DEPLOYMENT_NAME"]

PRIMARY_PRIVATE_KEY = os.environ["PRIMARY_PRIVATE_KEY_FIELD"]
PRIMARY_PUBLIC_KEY = os.environ["PRIMARY_PUBLIC_KEY_FIELD"]
SECONDARY_PRIVATE_KEY = os.environ["SECONDARY_PRIVATE_KEY_FIELD"]
SECONDARY_PUBLIC_KEY = os.environ["SECONDARY_PUBLIC_KEY_FIELD"]
KEY_PAIR_VERSION = os.environ["KEY_PAIR_VERSION_FIELD"]


def generate_key_pair():
    private_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption(),
    )
    public_pem = private_key.public_key().public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return private_pem, public_pem


def main():
    config.load_incluster_config()
    v1 = client.CoreV1Api()
    apps_v1 = client.AppsV1Api()

    try:
        secret = v1.read_namespaced_secret(SECRET_NAME, NAMESPACE)
        data = secret.data or {}
        current_primary_private = data.get(PRIMARY_PRIVATE_KEY)
        current_primary_public = data.get(PRIMARY_PUBLIC_KEY)
    except client.exceptions.ApiException as e:
        if e.status != 404:
            raise
        secret = client.V1Secret(
            metadata=client.V1ObjectMeta(name=SECRET_NAME, namespace=NAMESPACE),
            data={},
        )
        data = {}
        current_primary_private = None
        current_primary_public = None

    new_private_pem, new_public_pem = generate_key_pair()
    new_version = base64.b64encode(str(uuid.uuid4()).encode()).decode()

    if not current_primary_private or not current_primary_public:
        print("No existing keys found, generating both primary and secondary key pairs")
        secondary_private_pem, secondary_public_pem = generate_key_pair()
        secret.data[SECONDARY_PRIVATE_KEY] = base64.b64encode(secondary_private_pem).decode()
        secret.data[SECONDARY_PUBLIC_KEY] = base64.b64encode(secondary_public_pem).decode()
    else:
        secret.data[SECONDARY_PRIVATE_KEY] = current_primary_private
        secret.data[SECONDARY_PUBLIC_KEY] = current_primary_public

    secret.data[PRIMARY_PRIVATE_KEY] = base64.b64encode(new_private_pem).decode()
    secret.data[PRIMARY_PUBLIC_KEY] = base64.b64encode(new_public_pem).decode()
    secret.data[KEY_PAIR_VERSION] = new_version

    if data:
        v1.patch_namespaced_secret(SECRET_NAME, NAMESPACE, secret)
    else:
        v1.create_namespaced_secret(NAMESPACE, secret)
    print(f"Secret {SECRET_NAME} updated successfully")

    if os.environ.get("REDEPLOY", "true").lower() != "false":
        deployment = apps_v1.read_namespaced_deployment(DEPLOYMENT_NAME, NAMESPACE)
        annotations = deployment.spec.template.metadata.annotations or {}
        annotations["kubectl.kubernetes.io/restartedAt"] = datetime.datetime.utcnow().isoformat()
        deployment.spec.template.metadata.annotations = annotations
        apps_v1.patch_namespaced_deployment(DEPLOYMENT_NAME, NAMESPACE, deployment)
        print(f"Deployment {DEPLOYMENT_NAME} restarted successfully")
    else:
        print(f"Skipping redeploy of {DEPLOYMENT_NAME} (REDEPLOY=false)")


if __name__ == "__main__":
    main()
