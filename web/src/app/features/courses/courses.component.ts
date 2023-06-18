import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Course } from '../../api/models/course';
import { ISampleCoach } from '../../api/admin/admin-api.service';
import { ApiService } from '../../api/api.service';

@Component({
  selector: 'anms-courses',
  templateUrl: './courses.component.html',
  styleUrls: ['./courses.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CoursesComponent implements OnInit {

  public courses: Course[] = []

  constructor(
    private api: ApiService,
  ) {
  }

  ngOnInit(): void {
    this.api.getCourses().subscribe(v => this.courses = v)
  }

  // 显示助教信息
  public displayAssistantCoaches(coaches: ISampleCoach[]): string {
    return coaches.map(v => v.name).join('、')
  }
}
