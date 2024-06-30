import {Component, NgZone, OnDestroy, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';

interface Post {
    username: string;
    content: string;
    createdAt: string;
}

@Component({
    selector: 'app-timeline',
    templateUrl: './timeline.component.html',
    styleUrls: ['./timeline.component.css']
})
export class TimelineComponent implements OnInit, OnDestroy {
    username: string = 'Benutzername';
    newPostContent: string = '';
    timeline: Post[] = [];
    server: string = '';
    creatorName: string = '';
    searchQuery: string = '';
    searchResults: Post[] = [];
    private eventSource: EventSource | null = null;

    constructor(private http: HttpClient, private ngZone: NgZone) {
    }

    ngOnInit(): void {
        // this.subscribeToPosts();
    }

    ngOnDestroy() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }

    subscribeToPosts() {
        console.log("subscribeToPosts");
        this.eventSource = new EventSource('/api/timeline');
        this.eventSource.onmessage = (event) => {
            console.log(event);
            this.ngZone.run(() => {
                const post = JSON.parse(event.data);
                this.timeline.push(post);
            });
        };
    }

    createPost(): void {
        let content = this.newPostContent.trim();
        if (content !== '') {
            this.http.post('/api/post', { content: content }).subscribe(() => {
                this.newPostContent = '';
            });
        }
    }

    searchPosts(): void {
        let query = this.searchQuery.trim()
        if (query !== '') {
            this.searchResults = [];
            this.http.get<Post[]>(`/api/search?q=${query}`).subscribe({
                next: (results) => {
                    // this.ngZone.run(() => {
                        this.searchResults = results;
                    // });
                },
                error: (err) => {
                    console.log('Error!', err);
                }
            });
        }
    }

}
