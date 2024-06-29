import { Injectable } from '@angular/core';
import { Client, Message } from '@stomp/stompjs';
import { Observable, Subject } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class ChatService {
    private stompClient: Client;
    private messages: string[] = [];
    private messagesSubject: Subject<string[]> = new Subject<string[]>();
    public messages$: Observable<string[]> = this.messagesSubject.asObservable();

    constructor() {
        // Den WebSocket-Pfad konstruieren, der relativ zur aktuellen URL ist
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        const port = window.location.port ? ':' + window.location.port : '';

        this.stompClient = new Client({
            brokerURL: `${protocol}//${host}${port}/ws`,
            reconnectDelay: 5000,
            debug: (str: string) => {
                console.log(str);
            }
        });
        this.stompClient.onConnect = () => {
            this.stompClient.subscribe('/topic/messages', (message: Message) => {
                const lineUpdate = JSON.parse(message.body);
                const line = lineUpdate.line;
                const textMessage = lineUpdate.message;

                this.messages[line] = textMessage;

                this.messagesSubject.next(this.messages);
            });
        };
        this.stompClient.activate();
    }

    sendPartialMessage(message: string) {
        this.stompClient.publish({ destination: '/app/send-partial-message', body: message });
    }

    sendMessage(message: string) {
        this.stompClient.publish({ destination: '/app/send-message', body: message });
    }
}
