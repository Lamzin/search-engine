import { Component } from '@angular/core';

import { HttpClient } from '@angular/common/http';
import { PageEvent } from '@angular/material';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';
  query = '';
  results = [];

  constructor(private http: HttpClient) {}

  sendSearchRequest() {
    console.log('http://localhost:8080/search?query=' + this.query);
    this.http.get('http://localhost:8080/search?query=' + this.query)
      .subscribe(
        (results: string[]) => this.results = results,
        error => {
          console.log(error);
          this.results = [];
        }
      );
  }

  onPageEvent(event: PageEvent) {
    console.log(event);
  }

}
