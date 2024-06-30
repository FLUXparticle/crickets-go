import {Injectable, NgZone} from '@angular/core';
import {Observable, Subject} from 'rxjs';

export interface Post {
    username: string;
    content: string;
    createdAt: string;
}

@Injectable({
    providedIn: 'root'
})
export class ChatService {
    private ws: WebSocket;
    private messages: Post[] = [];
    private messagesSubject: Subject<Post[]> = new Subject<Post[]>();
    public messages$: Observable<Post[]> = this.messagesSubject.asObservable();

    constructor(private ngZone: NgZone) {
        // Den WebSocket-Pfad konstruieren, der relativ zur aktuellen URL ist
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        const port = window.location.port ? ':' + window.location.port : '';

        this.ws = new WebSocket(`${protocol}//${host}${port}/api/chatWS`);

        this.ws.onmessage = (event) => {
            this.ngZone.run(() => {
                const message = JSON.parse(event.data);
                this.messages.push(message);
                this.messagesSubject.next(this.messages);
            });
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onclose = () => {
            console.log('WebSocket connection closed');
        };
    }

    sendMessage(message: string) {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({ content: message }));
        } else {
            console.error('WebSocket is not open');
        }
    }
}
