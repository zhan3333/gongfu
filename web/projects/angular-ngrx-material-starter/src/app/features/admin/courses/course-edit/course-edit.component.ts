import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { AdminApiService, ISampleCoach } from '../../../../api/admin/admin-api.service';
import { NotificationService } from '../../../../core/notifications/notification.service';
import { ActivatedRoute } from '@angular/router';
import { FormControl, FormGroup } from '@angular/forms';
import { School } from '../../../../api/models/school';
import * as moment from 'moment';
import { Course } from '../../../../api/models/course';

@Component({
  selector: 'anms-course-edit',
  templateUrl: './course-edit.component.html',
  styleUrls: ['./course-edit.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CourseEditComponent implements OnInit {
  public course: Course | undefined;
  public id = 0

  // 教练数组
  public coaches: ISampleCoach[] = []
  public loading = false
  public schools: School[] = []

  public form = new FormGroup({
    schoolId: new FormControl(null),
    startDate: new FormControl(null),
    startTime: new FormControl(null),
    managerId: new FormControl(null),
    coachId: new FormControl(null),
    assistantCoachIds: new FormControl([]),
    summary: new FormControl(''),
    images: new FormControl([]),
    checkInAt: new FormControl(''),
    checkOutAt: new FormControl(''),
    checkInBy: new FormControl(''),
    checkOutBy: new FormControl(''),
  })

  constructor(
    private adminApi: AdminApiService,
    private notification: NotificationService,
    public route: ActivatedRoute,
  ) {
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.id = parseInt(params['id'], 10)
      if (this.id === 0) {
        this.notification.error('invalid request params')
        return
      }
      this.adminApi.getCoaches().subscribe(data => this.coaches = data)
      this.adminApi.getSchools().subscribe(data => this.schools = data)
      this.refreshCourse(this.id)
    })
  }

  refreshCourse(id: number) {
    this.adminApi.getCourse(id).subscribe(course => {
      this.course = course
      console.log('course', this.course)
      this.form.patchValue({
        schoolId: this.course.school.id,
        startDate: moment(this.course.startDate, 'YYYY/MM/DD').toDate(),
        startTime: this.course.startTime,
        managerId: this.course.manager.id,
        coachId: this.course.coach?.id,
        assistantCoachIds: this.course.assistantCoaches.map(item => item.id),
        summary: this.course.summary,
        images: this.course.images || [],
        checkInAt: this.course.checkInAt,
        checkOutAt: this.course.checkOutAt,
        checkInBy: this.course.checkInByUser?.id,
        checkOutBy: this.course.checkOutByUser?.id,
      })
    })
  }

  onSave() {
    if (!this.course?.id) {
      return
    }
    this.loading = true
    const data = this.form.value
    console.log(this.form.value)
    this.adminApi.updateCourse(this.course.id, {
      startDate: moment(data['startDate']).format('YYYY/MM/DD'),
      startTime: data['startTime'],
      schoolId: data['schoolId'],
      managerId: data['managerId'],
      coachId: data['coachId'],
      assistantCoachIds: data['assistantCoachIds'],
    }).subscribe(() => {
      this.notification.success('保存成功')
      this.refreshCourse(this.id)
    }, err => {
      this.notification.error('保存失败')
      console.error('保存失败', err)
    }, () => this.loading = false)
  }
}
