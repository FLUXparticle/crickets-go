import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Component({
    selector: 'app-profile',
    templateUrl: './profile.component.html',
    styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
    success: string = '';
    error: string = '';
    subscriberCount: number = 0;
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
                    this.error = `Error fetching profile data: ${err}`;
                }
            });
    }

    onSubmit(): void {
        this.http.post<{ success: string, error: string }>('/api/subscribe', this.subscribeForm)
            .subscribe({
                next: (data) => {
                    this.success = data.success || '';
                    this.error = data.error || '';
                },
                error: (err) => {
                    this.success = '';
                    this.error = `Error subscribing: ${err}`;
                }
            });
    }
}
