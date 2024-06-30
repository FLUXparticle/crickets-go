import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Component({
    selector: 'app-profile',
    templateUrl: './profile.component.html',
    styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
    subscriberCount: number = 0;
    successes: string[] = [];
    errors: string[] = [];
    subscribeForm = {
        server: '',
        creatorName: ''
    };

    constructor(private http: HttpClient) {}

    ngOnInit(): void {
        this.getProfile();
    }

    getProfile(): void {
        this.http.get<{ subscriberCount: number }>('/api/profile')
            .subscribe({
                next: (data) => {
                    this.subscriberCount = data.subscriberCount;
                },
                error: (err) => {
                    console.error('Error fetching profile data:', err);
                }
            });
    }

    onSubmit(): void {
        this.http.post<{ successes: string[], errors: string[] }>('/api/subscribe', this.subscribeForm)
            .subscribe({
                next: (data) => {
                    this.successes = data.successes || [];
                    this.errors = data.errors || [];
                    this.getProfile(); // Refresh subscriber count after subscribing
                },
                error: (err) => {
                    console.error('Error subscribing:', err);
                    this.errors = ['An error occurred while subscribing. Please try again.'];
                }
            });
    }
}
