import { Component, NgZone, OnDestroy, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface Post {
    creator: { username: string };
    content: string;
}

@Component({
    selector: 'app-timeline',
    templateUrl: './timeline.component.html',
    styleUrls: ['./timeline.component.css']
})
export class TimelineComponent implements OnInit, OnDestroy {
    username: string = 'Benutzername';
    newPostContent: string = '';
    posts: Post[] = [];
    server: string = '';
    creatorName: string = '';
    searchQuery: string = '';
    searchResults: Post[] = [];
    private eventSource: EventSource | null = null;
    private searchEventSource: EventSource | null = null;

    constructor(private http: HttpClient, private ngZone: NgZone) {
    }

    ngOnInit(): void {
        this.http.get<{ username: string }>('/api/username').subscribe(response => {
            this.username = response.username;
        });
        this.subscribeToPosts();
    }

    ngOnDestroy() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }

    subscribeToPosts() {
        console.log("subscribeToPosts");
        this.eventSource = new EventSource('/api/posts');
        this.eventSource.onmessage = (event) => {
            console.log(event);
            this.ngZone.run(() => {
                const post = JSON.parse(event.data);
                this.posts.push(post);
            });
        };
    }

    createPost(): void {
        if (this.newPostContent.trim() !== '') {
            this.http.post('/api/post', { content: this.newPostContent }).subscribe(() => {
                this.newPostContent = '';
            });
        }
    }

    searchPosts(): void {
        if (this.searchQuery.trim() !== '') {
            if (this.searchEventSource) {
                this.searchEventSource.close();
            }
            this.searchResults = [];
            this.searchEventSource = new EventSource(`/api/search?query=${this.searchQuery}`);
            this.searchEventSource.onmessage = (event) => {
                this.ngZone.run(() => {
                    const post = JSON.parse(event.data);
                    this.searchResults.push(post);
                });
            };
            this.searchEventSource.onerror = (event) => {
                console.log("Error!")
                this.searchEventSource?.close();
                this.searchEventSource = null;
            };
        }
    }

}
