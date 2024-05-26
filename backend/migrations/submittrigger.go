package migrations

import (
	"gorm.io/gorm"
)

type SubmitTriggerMigration struct{}

func CreateTrigger(db *gorm.DB) error {
	// 触发器，当 submissions 表的字段更新为 4 （Accepted) 的时候会触发
	createFunctionSQL := `
	CREATE OR REPLACE FUNCTION "public"."update_contest_rank"()
	  RETURNS "pg_catalog"."trigger" AS $BODY$
			DECLARE
				isaccepted BOOLEAN;
				contest_start_at TIMESTAMP;
				penalty_time INT8;
				contest INT8;
				penalty_c INT8;
			BEGIN
				IF tg_op = 'UPDATE' THEN
					-- 检查 NEW.id 是否是竞赛 id
					SELECT contest_id INTO contest
					FROM public.contest_submits
					WHERE submit_id = NEW.id;
					-- 找到当前的逻辑
					IF NOT FOUND THEN
						RETURN NEW;
					END IF;
					
					-- 检查该用户之前是否有 AC 提交
					SELECT accepted, penalty_count INTO isaccepted, penalty_c
					FROM public.contest_ranks
					WHERE contest_ranks.problem_id = NEW.problem_id
					AND contest_ranks.user_id = NEW.user_id
					AND contest_ranks.contest_id = contest;
					
					IF NEW.status = 4 THEN
						-- 计算罚时
						SELECT start_at INTO contest_start_at
						FROM public.contests
						WHERE id = contest;
						
						penalty_time := EXTRACT(EPOCH FROM (NEW.submit_time - contest_start_at)) + penalty_c * 20 * 60;
						-- 设置状态和罚时
						IF isaccepted = FALSE THEN
							UPDATE public.contest_ranks
							SET accepted = TRUE,
									penalty = penalty_time
							WHERE contest_ranks.problem_id = NEW.problem_id
							AND contest_ranks.user_id = NEW.user_id
							AND contest_ranks.contest_id = contest;
						END IF;
					ELSIF NEW.status >= 5 AND NEW.status <= 10 THEN
							-- 更新罚时次数		
					RAISE NOTICE 'An insert operation occurred on %', TG_TABLE_NAME;
						IF isaccepted = FALSE THEN
							UPDATE public.contest_ranks
							SET penalty_count = penalty_count + 1
							WHERE contest_ranks.problem_id = NEW.problem_id
								AND contest_ranks.user_id = NEW.user_id
								AND contest_ranks.contest_id = contest;
						END IF;
				END IF;		
			END IF;
			RETURN NEW;
		END $BODY$
	  LANGUAGE plpgsql VOLATILE
	  COST 100
    `
	// 执行创建触发器函数的SQL语句
	if err := db.Exec(createFunctionSQL).Error; err != nil {
		return err
	}
	// 检查触发器是否已存在
	var count int
	err := db.Raw("SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name = 'update_contest_rank_trigger' AND event_object_table = 'submissions'").Scan(&count).Error
	if err != nil {
		return err
	}

	// 创建触发器的SQL语句
	createTriggerSQL := `
    CREATE TRIGGER update_contest_rank_trigger
    AFTER INSERT OR UPDATE ON submissions
    FOR EACH ROW
    EXECUTE FUNCTION update_contest_rank();
    `
	if count == 0 {
		if err := db.Exec(createTriggerSQL).Error; err != nil {
			// 注意：如果触发器因其他原因（如权限问题）而创建失败，这里也会返回错误
			return err
		}
	}
	return nil
}

func (m *SubmitTriggerMigration) Rollback(db *gorm.DB) error {
	return db.Exec(`
        DROP FUNCTION IF EXISTS update_contest_rank();
    `).Error
}
