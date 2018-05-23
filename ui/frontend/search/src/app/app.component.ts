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
  resultsCount = 0;

  constructor(private http: HttpClient) {}

  sendSearchRequest(page: number = 0, pageSize: number=10) {
    var url = `http://localhost:8080/search?query=${this.query};skip=${page*pageSize};limit=${pageSize}`;
    console.log(url);
    this.http.get(url)
      .subscribe(
        (response: any) => {
          this.results = response.results;
          this.resultsCount = response.results_count;
        },
        error => {
          console.log(error);
          this.results = [];
        }
      );
  }

  onPageEvent(event: PageEvent) {
    console.log(event);
    this.sendSearchRequest(event.pageIndex, event.pageSize);
  }

}
