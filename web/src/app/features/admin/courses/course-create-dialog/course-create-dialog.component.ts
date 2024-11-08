import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormsModule, NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AdminApiService, ISampleCoach } from '../../../../api/admin/admin-api.service';
import { NotificationService } from '../../../../core/notifications/notification.service';
import { Router } from '@angular/router';
import { School } from '../../../../api/models/school';
import * as moment from 'moment';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatInputModule } from '@angular/material/input';
import { MatOptionModule } from '@angular/material/core';
import { NgFor } from '@angular/common';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'anms-course-create-dialog',
  templateUrl: './course-create-dialog.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [MatCardModule, FormsModule, ReactiveFormsModule, MatFormFieldModule, MatSelectModule, NgFor, MatOptionModule, MatInputModule, MatDatepickerModule, MatButtonModule, MatDialogModule]
})
export class CourseCreateDialogComponent implements OnInit {
  public form = this.fb.group({
    schoolId: [0, Validators.required],
    startDate: ['', [Validators.required]],
    startTime: ['', [Validators.required]],
    managerId: [0, [Validators.required]],
    coachId: [0, [Validators.required]],
    assistantCoachIds: [[] as Array<number>]
  });

  // 教练数组
  public coaches: ISampleCoach[] = [];
  public loading = false;
  public schools: School[] = [];

  constructor(
    public dialogRef: MatDialogRef<CourseCreateDialogComponent>,
    private adminApi: AdminApiService,
    private notification: NotificationService,
    private router: Router,
    private fb: NonNullableFormBuilder
  ) {
  }

  ngOnInit(): void {
    this.adminApi.getCoaches().subscribe((data) => (this.coaches = data));
    this.adminApi.getSchools().subscribe((data) => (this.schools = data));
  }

  onSave() {
    const data = this.form.value;
    this.loading = true;
    this.adminApi
      .createCourse({
        assistantCoachIds: this.form.value.assistantCoachIds || [],
        coachId: this.form.value.coachId || 0,
        managerId: this.form.value.managerId || 0,
        startDate: moment(this.form.value.startDate).format('YYYY/MM/DD'),
        startTime: this.form.value.startTime || '',
        schoolId: this.form.value.schoolId || 0,
      })
      .subscribe(
        () => {
          this.notification.success('创建成功');
          this.dialogRef.close(true)
        },
        (err) => {
          this.notification.error('创建失败');
          console.error('create course failed', err);
        },
        () => (this.loading = false)
      );
  }
}
