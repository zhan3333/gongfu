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
import { StudyRecord } from '../../../api/models/study-record';
import { NgIf } from '@angular/common';
import * as moment from 'moment/moment';

@Component({
  selector: 'anms-teaching-record-dialog',
  template: `
    <h1 mat-dialog-title>{{ data ? '编辑学习记录' : '新增学习记录' }}</h1>
    <div mat-dialog-content [formGroup]="form">
      <mat-form-field>
        <mat-label>学习日期</mat-label>
        <input
          matInput
          [matDatepicker]="picker"
          formControlName="date"
          [disabled]="true"
        />
        <mat-hint>MM/DD/YYYY</mat-hint>
        <mat-datepicker-toggle
          matIconSuffix
          [for]="picker"
        ></mat-datepicker-toggle>
        <mat-datepicker #picker></mat-datepicker>
      </mat-form-field>
      <mat-form-field>
        <mat-label>学习内容</mat-label>
        <input matInput formControlName="content" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions>
      <button mat-button (click)="onNoClick()">取消</button>
      <button
        mat-button
        (click)="onClickDelete()"
        *ngIf="form.value.id"
        color="warn"
      >
        删除
      </button>
      <button mat-button (click)="onClickSubmit()" color="primary">保存</button>
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
export class StudyRecordDialog {
  public form = this.fb.group({
    id: [0],
    date: ['', Validators.required],
    content: ['', Validators.required]
  });

  constructor(
    public dialogRef: MatDialogRef<StudyRecordDialog>,
    private fb: FormBuilder,
    @Inject(MAT_DIALOG_DATA) public data?: StudyRecord
  ) {
    if (data) {
      this.form.patchValue({
        id: data.id,
        date: moment(data.date, 'YYYY/MM/DD').format(),
        content: data.content
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
        content: this.form.value.content
      }
    });
  }

  onClickDelete(): void {
    this.dialogRef.close({ isDelete: true });
  }
}
