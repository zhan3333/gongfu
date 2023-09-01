import { Component, Inject } from '@angular/core';
import {
  MAT_DIALOG_DATA,
  MatDialogModule,
  MatDialogRef
} from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import {
  FormBuilder,
  FormsModule,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { TeachingRecord } from '../../../api/models/teaching-record';
import { NgIf } from '@angular/common';
import * as moment from 'moment';

@Component({
  selector: 'anms-teaching-record-dialog',
  template: `
    <h1 mat-dialog-title>{{ data ? '编辑授课记录' : '新增授课记录' }}</h1>
    <div mat-dialog-content [formGroup]="form">
      <mat-form-field>
        <mat-label>授课日期</mat-label>
        <input matInput [matDatepicker]="picker" formControlName="date" />
        <mat-hint>MM/DD/YYYY</mat-hint>
        <mat-datepicker-toggle
          matIconSuffix
          [for]="picker"
        ></mat-datepicker-toggle>
        <mat-datepicker #picker [disabled]="false"></mat-datepicker>
      </mat-form-field>
      <mat-form-field>
        <mat-label>授课内容</mat-label>
        <input matInput formControlName="address" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions>
      <button mat-button (click)="onNoClick()">取消</button>
      <button
        mat-raised-button
        color="warn"
        (click)="onClickDelete()"
        *ngIf="form.value.id"
      >
        删除
      </button>
      <button
        mat-raised-button
        (click)="onClickSubmit()"
        [disabled]="form.invalid"
        color="primary"
      >
        保存
      </button>
    </div>
  `,
  standalone: true,
  imports: [
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    MatButtonModule,
    ReactiveFormsModule,
    MatDatepickerModule,
    NgIf
  ]
})
export class TeachingRecordDialog {
  public form = this.fb.group({
    id: [0],
    date: ['', Validators.required],
    address: ['', Validators.required]
  });

  constructor(
    public dialogRef: MatDialogRef<TeachingRecordDialog>,
    private fb: FormBuilder,
    @Inject(MAT_DIALOG_DATA) public data?: TeachingRecord
  ) {
    if (data) {
      this.form.patchValue({
        id: data.id,
        date: moment(data.date, 'YYYY/MM/DD').format(),
        address: data.address
      });
    }
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  onClickSubmit(): void {
    this.dialogRef.close({
      data: {
        id: this.form.value.id,
        date: moment(this.form.value.date).format('YYYY/MM/DD'),
        address: this.form.value.address
      }
    });
  }

  onClickDelete(): void {
    this.dialogRef.close({ isDelete: true });
  }
}
