# Iniciando a configuração do servidor

## Configurando a Internet via USB Tethering

Com o seu smartphone conectado via cabo USB, siga os passos abaixo:

1. Ative o USB Tethering no Smartphone:
   No Android (ou iOS), acesse:
   - Configurações → Rede e Internet → Tethering e Hotspot → Ative Tethering USB.
2. Execute o comando para identificar a nova interface (ex: usb0 ou enx...):
   ```bash
   ip a
   ```
3. Adicione o IP manualmente para sua interface:
   ```bash
   sudo ip addr add 172.20.10.2/24 dev enx5e52842e9eb2
   ```
   **Dica:** Substitua `enx5e52842e9eb2` pelo nome da interface identificada no passo anterior.
4. Ative a Interface executando o comando:
   ```bash
   sudo ip link set enx5e52842e9eb2 up
   ```
   Em seguida, adicione a rota padrão usando o comando:
   ```bash
   sudo ip route add default via 172.20.10.1 dev enx5e52842e9eb2
   ```
5. Confirme o IP obtido:
   ```bash
   ip a show enx5e52842e9eb2
   ```
6. Configure o DNS adicionando o endereço do servidor DNS primário do Google:
   ```bash
   echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
   ```
   Verifique se o arquivo foi atualizado
   ```bash
   cat /etc/resolv.conf
   ```
   Deverá aparecer algo escrito: `nameserver 8.8.8.8`
7. Teste a conexão usando os seguintes comandos no seu terminal:
   ```bash
   ping -c 4 172.20.10.11
   ping -c 4 8.8.8.8
   ping -4 google.com
   ```

## Tornando a Configuração Permanente

Agora, vamos instalar os pacotes necessários e configurar o NetworkManager para que a conexão seja persistente após as reinicializações do sistema.

1. Atualize a lista de pacotes e instale o **NetworkManager**, **wpasupplicant** e **linux-firmware**:
   ```bash
   sudo apt update
   sudo apt install network-manager wpasupplicant linux-firmware -y
   ```
2. Configure o Netplan para o NetworkManager editando o arquivo de configuração do netplan:
   ```bash
   sudo nano /etc/netplan/01-netcfg.yaml
   ```
   Altere o conteúdo do arquivo para o descrito abaixo e salve o arquivo usando (Ctrl + O, Enter, Ctrl + X).:
   ```yaml
   network:
     version: 2
     renderer: NetworkManager
   ```
   É imporante definir as permissões do arquivo para garantir acesso apenas pelo root:
   ```bash
   sudo chmod 600 /etc/netplan/01-netcfg.yaml
   ```
3. Aplique as alterações do netplan:
   ```bash
   sudo netplan apply
   ```
4. Conectando-se à Rede Wi-Fi:
   Você pode conectar-se à rede Wi-Fi utilizando a linha de comando ou a interface gráfica do NetworkManager.

   - Via CLI:

   ```bash
   nmcli device wifi connect "<nome-do-wifi>" password "<senha-do-wifi>"
   ```

   Ou via Interface Gráfica:

   ```bash
   nmtui
   ```

5. Remova o cabo USB do seu computador e rode o comando para checar se a rede está funcionando:
   ```bash
   ping google.com
   ```
6. Por fim, reinicie seu sistema e cheque se internet continua funcionando:
   ```bash
   sudo reboot
   ping google.com
   ```

## Configurando o SSH

Agora que a internet está funcionando, vamos configurar o SSH para acessar o servidor remotamente.

1. Instale o OpenSSH Server:
   ```bash
   sudo apt install openssh-server -y
   ```
2. Verifique se o serviço SSH está ativo:
   ```bash
    sudo systemctl status ssh
   ```
   Caso o serviço não esteja ativo, inicie-o: `sudo systemctl start ssh`
3. Para garantir que o SSH inicie automaticamente após a reinicialização, execute:
   ```bash
    sudo systemctl enable ssh
   ```
4. Verifique o IP do servidor:
   ```bash
   hostname -I
   ```
5. Conecte-se ao servidor via SSH a partir de outro computador:
   ```bash
   ssh <usuario>@<ip-do-servidor>
   ```

## Adicionando chave SSH no servidor

1. Gere uma chave SSH no seu computador local:
   ```bash
   ssh-keygen -t rsa -b 4096
   ```
2. Copie a chave pública para o servidor:
   ```bash
   ssh-copy-id <usuario>@<ip-do-servidor>
   # Ou o comando abaixo caso você esteja utilizando windows
   type $env:USERPROFILE\.ssh\id_rsa.pub | ssh <usuario>@<ip-do-servidor> "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
   ```
3. Verifique se a chave foi adicionada corretamente:
   ```bash
   ssh <usuario>@<ip-do-servidor>
   ```
4. Para aumentar a segurança, desative o login por senha editando o arquivo de configuração do SSH:
   ```bash
   sudo nano /etc/ssh/sshd_config
   ```
   Altere a linha `PasswordAuthentication yes` para `PasswordAuthentication no` e salve o arquivo.
5. Reinicie o serviço SSH para aplicar as alterações:
   ```bash
   sudo systemctl restart ssh
   ```
6. Verifique se o SSH está funcionando corretamente sem senha:
   ```bash
   ssh <usuario>@<ip-do-servidor>
   ```

## Acessando o servido através da internet pública com SSH via Cloudflared Tunnel

1. Crie uma conta no Cloudflare e adicione seu domínio.
2. Entre no painel do Cloudflare, vá em Zero trust, depois em Networks e clique em Tunnels.
3. Clique em "Create a tunnel" e siga as instruções para criar um túnel.
4. Escolha "Cloudflared" como seu tipo de tunnel e clique em "Next".
5. Dê um nome para o seu túnel e clique em "Save tunnel".
6. Agora, é necessário baixar o cloudflared. Você pode selecionar o seu sistema operacional na tela de "Install and run connectors" e rodar os comandos sugeridos que serão parecidos como:

   ```bash
   sudo mkdir -p --mode=0755 /usr/share/keyrings
   curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null

   echo 'deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared any main' | sudo tee /etc/apt/sources.list.d/cloudflared.list

   sudo apt-get update && sudo apt-get install cloudflared
   ```

7. Agora você precisará linkar o domínio que cadastrou no inicio com algum serviço rodando no seu ubuntu server. Recomendo que você coloque o seu ssh para que seja possível acessar atráves da internet pública.
8. Faça login no cloudflared rodando o comando:
   ```bash
   cloudflared login
   ```
   Isso gerará um link que você precisa abrir e selecionar o domínio que você cadastrou.
9. Agora você pode iniciar o seu tunneling usando o comando:

```bash
 cloudflared tunnel run <nome-do-tunnel>
```

10. Adicione um subdomínio para o seu túnel usando o painel do Cloudflare:
    **Dica:** O subdomínio deve ser algo como `ssh.dominio.com` e o domínio deve ser o mesmo que você cadastrou no Cloudflare.
11. Agora você pode acessar o seu servidor através do subdomínio que você cadastrou.

```bash
ssh -o ProxyCommand="cloudflared access ssh --hostname %h" ruan-homelab@ssh.whoam.site
```

Lembre-se: para acessar o ssh é necessário que o client também tenha o cloudflared instalado.

## Instalando Docker

1. Atualize o sistema

```bash
sudo apt update && sudo apt upgrade -y
```

2. Instale dependências

```bash
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
```

3. Adicione a chave do Docker

```bash
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

4. Instale o Docker

```bash
sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io
```

5. Adicione seu usuário ao grupo Docker

```bash
sudo usermod -aG docker $USER
newgrp docker
```

## Instalando o K3s

1. Instale o K3s

```bash
curl -sfL https://get.k3s.io | sh -s - --docker
```

3. Configure o acesso ao cluster:

```bash
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config
export KUBECONFIG=~/.kube/config
```

2. Verifique a instalação:

```bash
kubectl get nodes  # Deve mostrar "Ready"
kubectl get pods -A  # Verifique o Traefik (já vem pré-instalado)
```

## Instalando o ArgoCD:

1. Instale o helm usando o comando:

   ```bash
   sudo snap install helm --classic
   ```

2. Instale o ArgoCD:

```bash
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

3.

## Instalando o Hashicorp Vault:

1. Adicione o repositório oficial da Hashicorp:

```bash
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
sudo tee /etc/apt/sources.list.d/hashicorp.list

sudo apt update && sudo apt install vault -y
```

2. Configure o Vault com Systemd e Raft (sem TLS)

   1. Crie o arquivo de configuração:

   ```bash
   sudo mkdir -p /opt/vault/data
   sudo chown -R vault:vault /opt/vault
   sudo tee /etc/vault.d/vault.hcl > /dev/null <<EOF
   ui = true
   disable_mlock = true

   storage "raft" {
      path    = "/opt/vault/data"
      node_id = "node1"
   }

   listener "tcp" {
      address     = "0.0.0.0:8200"
      tls_disable = 1
      cluster_addr = "http://<IP_SERVIDOR>:8201"
   }

   api_addr = "http://<IP_SERVIDOR>:8200"
   cluster_addr = "http://<IP_SERVIDOR>:8201"
   EOF
   ```

   2. Ative e inicie o Vault:

   ```bash
   sudo systemctl enable vault
   sudo systemctl start vault
   ```

3. Inicialize e desbloqueie o Vault:

   1. Exporte o vault_addr adicionado o texto abaixo no final do arquivo ~/.bashrc:

   ```bash
   nano ~/.bashrc
   ```

   Adicione o texto `export VAULT_ADDR='http://<IP-DO-SERVIDOR>:8200'` no final do arquivo

   2. Inicialize e guarde unseal keys e o root token com segurança:

   ```bash
   vault operator init
   ```

   3. Desbloqueie (unseal):

   ```bash
   vault operator unseal <unseal_key_1>
   vault operator unseal <unseal_key_2>
   vault operator unseal <unseal_key_3>
   ```

   4. login:

   ```bash
   vault login <root_token>
   ```

4. Habilitar KV v2 em secret/

```bash
vault secrets enable -path=secret kv-v2
vault secrets list
# deve listar secret/ kv-v2
```

5. Crie uma política e role no Vault

   1. Crie a política:

   ```bash
   vault policy write myapp-policy - <<EOF
   path "secret/data/myapp/*" {
   capabilities = ["read","create","update","delete","list"]
   }
   EOF
   ```

   2. Subir a policy:

   ```bash
   vault policy write myapp-policy myapp-policy.hcl
   ```

6. Configurar autenticação Kubernetes

   1. Obter dados do cluster:

   ```bash
   # CA cert
   kubectl config view --raw \
   -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' \
   | base64 -d > /tmp/ca.crt

   # API server
   KUBERNETES_HOST=$(kubectl config view --raw \
   -o jsonpath='{.clusters[0].cluster.server}')
   ```

   2. Crie o ServiceAccount e faça o bind da role:

   ```bash
   kubectl create sa vault-auth -n kube-system

   kubectl create clusterrolebinding vault-auth-delegator \
   --clusterrole=system:auth-delegator \
   --serviceaccount=kube-system:vault-auth
   ```

   3. Habilitar e configurar

   ```bash
   vault auth enable kubernetes

   vault write auth/kubernetes/config \
   token_reviewer_jwt="$(kubectl create token vault-auth -n kube-system)" \
   kubernetes_host="$KUBERNETES_HOST" \
   kubernetes_ca_cert=@/tmp/ca.crt
   ```

7. Criar role para injeção:

```bash
vault write auth/kubernetes/role/myapp \
  bound_service_account_names=vault-injector \
  bound_service_account_namespaces=default \
  policies=myapp-policy \
  ttl=24h
```

8. Instalar Vault Agent Injector no K3s:

```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

helm install vault hashicorp/vault \
  --namespace vault --create-namespace \
  --set "injector.enabled=true" \
  --set "server.enabled=false" \
  --set "externalVaultAddr=http://<IP_SERVIDOR>:8200"
```

9. Preparar a ServiceAccount no Kubernetes

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vault-injector
  namespace: default
EOF
```

10. Criar e testar o segredo
1. Grave no Vault:

```bash
vault kv put secret/mschat/config username="admin" password="s3cr3t"
```

2.  Leia para confirmar:

```bash
vault kv get secret/mschat/config
```

11. Deploy do seu app com anotações

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
      annotations:
        vault.hashicorp.com/agent-inject: "true"
        vault.hashicorp.com/role: "mschat"
        vault.hashicorp.com/agent-inject-secret-config.txt: "secret/data/mschat/config"
        vault.hashicorp.com/vault-addr: "http://192.168.15.40:8200"
    spec:
      serviceAccountName: vault-injector
      containers:
        - name: app
          image: busybox
          command: ["sleep", "3600"]
        - name: vault-agent
          image: hashicorp/vault:1.15.4
          volumeMounts:
            - name: vault-config
              mountPath: /home/vault
            - name: vault-config-secret
              mountPath: /vault/secret
          env:
            - name: VAULT_CONFIG
              value: /home/vault/config.json
      volumes:
        - name: vault-config
          emptyDir: {}
        - name: vault-config-secret
          secret:
            secretName: vault-agent-config
```

12. Verificar dentro do pod

```bash
kubectl exec -it deploy/myapp -- cat /vault/secrets/config.txt
# deve mostrar:
# username = admin
# password = s3cr3t
```
