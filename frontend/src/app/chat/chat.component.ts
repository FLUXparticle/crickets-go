import { Component, OnInit } from '@angular/core';
import { ChatService } from './chat.service';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent implements OnInit {
  newMessage: string = '';
  messages: string[] = [];

  constructor(private chatService: ChatService) { }

  ngOnInit(): void {
    this.chatService.messages$.subscribe((messages: string[]) => {
      this.messages = messages;
    });
  }

  sendPartialMessage() {
    this.chatService.sendPartialMessage(this.newMessage);
  }

  sendMessage() {
    this.chatService.sendMessage(this.newMessage);
    this.newMessage = '';
  }
}
