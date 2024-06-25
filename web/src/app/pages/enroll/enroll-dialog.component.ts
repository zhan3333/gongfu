import { Component, Inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-enroll-dialog',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatInputModule, ReactiveFormsModule, MatSelectModule, MatButtonModule, MatDialogModule],
  templateUrl: './enroll-dialog.component.html',
})
export class EnrollDialogComponent implements OnInit {
  public form = new FormGroup({
    username: new FormControl('', [Validators.required]),
    sex: new FormControl('', [Validators.required]),
    phone: new FormControl('', [Validators.required]),
  })

  constructor(
    public dialogRef: MatDialogRef<EnrollDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {
  }

  ngOnInit(): void {

  }

  confirm() {
    this.dialogRef.close(this.form.value);
  }
}
