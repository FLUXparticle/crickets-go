import {Component} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
    username = 'Unknown';

    constructor(private http: HttpClient) {}

    ngOnInit(): void {
        this.loadUsername();
    }

    loadUsername(): void {
        this.http.get<{ username: string }>('/api/username')
            // .subscribe(response => {})
            .subscribe({
                next: response => {
                    this.username = response.username;
                },
                error: error => {
                    console.error('Error loading username:', error);
                }
            });
    }
}
