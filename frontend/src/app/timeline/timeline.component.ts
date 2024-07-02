import {Component, NgZone, OnDestroy, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';

interface Post {
    creatorName: string;
    content: string;
    createdAt: string;
}

@Component({
    selector: 'app-timeline',
    templateUrl: './timeline.component.html',
    styleUrls: ['./timeline.component.css']
})
export class TimelineComponent implements OnInit, OnDestroy {
    errorPost: string = '';
    errorSearch: string = '';
    newPostContent: string = '';
    timeline: Post[] = [];
    server: string = '';
    creatorName: string = '';
    searchServer: string = '';
    searchQuery: string = '';
    searchResults: Post[] = [];
    private eventSource: EventSource | null = null;

    constructor(private http: HttpClient, private ngZone: NgZone) {
    }

    ngOnInit(): void {
        this.subscribeToTimeline();
    }

    ngOnDestroy() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }

    subscribeToTimeline() {
        this.eventSource = new EventSource('/api/timeline');
        this.eventSource.onmessage = (event) => {
            console.log(event);
            this.ngZone.run(() => {
                const post = JSON.parse(event.data);
                this.timeline.push(post);
            });
        };
        this.eventSource.onerror = (event) => {
            console.error("EventSource failed:", event);
            // Versuche die Verbindung nach einer kurzen Pause wiederherzustellen
            setTimeout(() => {
                this.subscribeToTimeline();
            }, 5000); // 5 Sekunden warten
        };
    }

    createPost(): void {
        let content = this.newPostContent.trim();
        if (content !== '') {
            this.http.post('/api/post', { content: content }).subscribe({
                next: () => {
                    this.newPostContent = '';
                },
                error: (err) => {
                    this.errorPost = `Error: ${err}`;
                }
            });
        }
    }

    searchPosts(): void {
        let server = this.searchServer.trim()
        let query = this.searchQuery.trim()
        if (query !== '') {
            this.searchResults = [];
            this.http.get<{searchResults:Post[],error:string}>(`/api/search?s=${server}&q=${query}`).subscribe({
                next: (data) => {
                    this.errorSearch = data.error || '';
                    this.searchResults = data.searchResults || [];
                },
                error: (err) => {
                    this.errorSearch = `Error: ${err}`;
                }
            });
        }
    }

}
