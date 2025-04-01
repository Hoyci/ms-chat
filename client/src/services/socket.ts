import { IMessage } from "@store/message";

class SocketService {
  private socket: WebSocket | null = null;

  connect(url: string) {
    this.socket = new WebSocket(url);
    this.socket.onopen = () => {
      console.log("WebSocket conectado");
    };
    this.socket.onmessage = (event) => {
      console.log("Mensagem recebida:", event.data);
    };
    this.socket.onclose = () => {
      console.log("WebSocket desconectado");
    };
    this.socket.onerror = (error) => {
      console.error("WebSocket erro:", error);
    };
  }

  sendMessage(message: IMessage) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.warn("Socket connection is not open");
    }
  }

  disconnect() {
    this.socket?.close();
  }
}

export default new SocketService();
