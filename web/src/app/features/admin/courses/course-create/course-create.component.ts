import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { UntypedFormControl, UntypedFormGroup, Validators, FormsModule, ReactiveFormsModule } from '@angular/forms';
import {
  AdminApiService,
  ISampleCoach
} from '../../../../api/admin/admin-api.service';
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

@Component({
    selector: 'anms-course-create',
    templateUrl: './course-create.component.html',
    styleUrls: ['./course-create.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true,
    imports: [MatCardModule, FormsModule, ReactiveFormsModule, MatFormFieldModule, MatSelectModule, NgFor, MatOptionModule, MatInputModule, MatDatepickerModule, MatButtonModule]
})
export class CourseCreateComponent implements OnInit {
  public form = new UntypedFormGroup({
    schoolId: new UntypedFormControl(null, [Validators.required]),
    startDate: new UntypedFormControl(null, [Validators.required]),
    startTime: new UntypedFormControl(null, [Validators.required]),
    managerId: new UntypedFormControl(null, [Validators.required]),
    coachId: new UntypedFormControl(null),
    assistantCoachIds: new UntypedFormControl([])
  });

  // 教练数组
  public coaches: ISampleCoach[] = [];
  public loading = false;
  public schools: School[] = [];

  constructor(
    private adminApi: AdminApiService,
    private notification: NotificationService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.adminApi.getCoaches().subscribe((data) => (this.coaches = data));
    this.adminApi.getSchools().subscribe((data) => (this.schools = data));
  }

  onSave() {
    const data = this.form.value;
    this.loading = true;
    this.adminApi
      .createCourse({
        assistantCoachIds: data['assistantCoachIds'],
        coachId: data['coachId'],
        managerId: data['managerId'],
        startDate: moment(data['startDate']).format('YYYY/MM/DD'),
        startTime: data['startTime'],
        schoolId: data['schoolId']
      })
      .subscribe(
        () => {
          this.notification.success('创建成功');
          this.router.navigate(['/admin/courses']);
        },
        (err) => {
          this.notification.error('创建失败');
          console.error('create course failed', err);
        },
        () => (this.loading = false)
      );
  }
}
