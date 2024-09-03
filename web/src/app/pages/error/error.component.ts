import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';

@Component({
  selector: 'app-error',
  standalone: true,
  imports: [CommonModule, MatCardModule],
  template: `
    <p>origin: {{ v }}</p>
    <mat-card *ngIf="error">
      <mat-card-header>
        <mat-card-title>{{ error.name }}</mat-card-title>
        <mat-card-subtitle>{{ error.message }}</mat-card-subtitle>
      </mat-card-header>
      <mat-card-content>
        <p>stack: {{ error.stack }}</p>
        <p>cause: {{ error.cause }}</p>
      </mat-card-content>
    </mat-card>`,
})
export class ErrorComponent implements OnInit {
  public v = ''
  public error: Error | undefined

  constructor() {
  }

  ngOnInit(): void {
    this.v = localStorage.getItem("error") || ''
    if (this.v) {
      this.error = JSON.parse(this.v);
    }
  }
}
