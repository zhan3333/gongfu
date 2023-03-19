import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { AdminApiService, ISampleCoach } from '../../../../api/admin/admin-api.service';
import { NotificationService } from '../../../../core/notifications/notification.service';
import { ActivatedRoute, Route, Router } from '@angular/router';

@Component({
  selector: 'anms-course-create',
  templateUrl: './course-create.component.html',
  styleUrls: ['./course-create.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CourseCreateComponent implements OnInit {
  public form = new FormGroup({
    name: new FormControl(''),
    address: new FormControl(''),
    schoolStartAt: new FormControl(''),
    coachId: new FormControl(null),
    assistantCoachIds: new FormControl([]),
  })

  // 教练数组
  public coaches: ISampleCoach[] = []
  public loading = false

  constructor(
    private adminApi: AdminApiService,
    private notification: NotificationService,
    private router: Router,
  ) {
  }

  ngOnInit(): void {
    this.adminApi.getCoaches().subscribe(data => this.coaches = data)
  }

  onSave() {
    console.log('save', this.form.value)
    const data = this.form.value
    this.loading = true
    this.adminApi.createCourse({
      address: data['address'],
      assistantCoachIds: data['assistantCoachIds'],
      coachId: data['coachId'],
      name: data['name'],
      schoolStartAt: new Date(data['schoolStartAt']).getTime() / 1000
    }).subscribe(
      () => {
        this.notification.success('创建成功')
        this.router.navigate(['/admin/courses'])
      },
      (err) => {
        this.notification.error('创建失败')
        console.error('create course failed', err)
      },
      () => this.loading = false)
  }
}
